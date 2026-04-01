package services

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/CodeEnthusiast09/country-currency-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CountryService struct {
	db              *gorm.DB
	externalService *ExternalService
	refreshMu       sync.Mutex
}

func NewCountryService(db *gorm.DB, externalService *ExternalService) *CountryService {
	return &CountryService{
		db:              db,
		externalService: externalService,
	}
}

func (s *CountryService) Refresh() (int, error) {
	s.refreshMu.Lock()

	defer s.refreshMu.Unlock()

	externalCountries, err := s.externalService.FetchCountries()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch countries: %w", err)
	}

	exchangeRates, err := s.externalService.FetchExchangeRates()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}

	now := time.Now()

	for _, ec := range externalCountries {
		country := models.Country{
			Name:            ec.Name,
			Capital:         ec.Capital,
			Region:          ec.Region,
			Population:      ec.Population,
			FlagURL:         ec.Flag,
			LastRefreshedAt: now,
		}

		if len(ec.Currencies) > 0 {
			code := ec.Currencies[0].Code
			country.CurrencyCode = &code

			if rate, ok := exchangeRates[code]; ok {
				country.ExchangeRate = &rate

				// estimated GDP = (population / exchange_rate) * random multiplier (1000-2000)
				multiplier := float64(rand.Intn(1001) + 1000)
				gdp := (float64(ec.Population) / rate) * multiplier
				country.EstimatedGDP = &gdp
			}
		}

		// upsert — update if name exists, insert if not
		// clause.OnConflict tells GORM: if name already exists, update all other columns
		result := s.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"capital", "region", "population",
				"currency_code", "exchange_rate", "estimated_gdp",
				"flag_url", "last_refreshed_at",
			}),
		}).Create(&country)

		if result.Error != nil {
			return 0, fmt.Errorf("failed to upsert country %s: %w", ec.Name, result.Error)
		}
	}

	return len(externalCountries), nil
}

func (s *CountryService) GetAll(region, currency, sort string) ([]models.Country, error) {
	var countries []models.Country

	query := s.db.Model(&models.Country{})

	if region != "" {
		query = query.Where("region = ?", region)
	}

	if currency != "" {
		query = query.Where("currency_code = ?", strings.ToUpper(currency))
	}

	switch sort {
	case "gdp_asc":
		query = query.Order("estimated_gdp ASC")
	case "gdp_desc":
		query = query.Order("estimated_gdp DESC")
	case "population_asc":
		query = query.Order("population ASC")
	case "population_desc":
		query = query.Order("population DESC")
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	}

	if err := query.Find(&countries).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}

	return countries, nil
}

func (s *CountryService) GetTopCountries(limit int) ([]models.Country, error) {
	var countries []models.Country
	if err := s.db.Order("estimated_gdp DESC").Limit(limit).Find(&countries).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch top countries: %w", err)
	}
	return countries, nil
}

func (s *CountryService) GetOne(name string) (*models.Country, error) {
	var country models.Country

	result := s.db.Where("LOWER(name) = LOWER(?)", name).First(&country)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("country not found")
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &country, nil
}

func (s *CountryService) Delete(name string) error {
	result := s.db.Where("LOWER(name) = LOWER(?)", name).Delete(&models.Country{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete country: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}

func (s *CountryService) GetStatus() (StatusResponse, error) {
	var total int64
	if err := s.db.Model(&models.Country{}).Count(&total).Error; err != nil {
		return StatusResponse{}, fmt.Errorf("failed to count countries: %w", err)
	}

	var latest models.Country
	if err := s.db.Order("last_refreshed_at DESC").First(&latest).Error; err != nil {
		return StatusResponse{Total: total}, nil
	}

	t := latest.LastRefreshedAt.Format(time.RFC3339)
	return StatusResponse{Total: total, LastRefreshedAt: &t}, nil
}
