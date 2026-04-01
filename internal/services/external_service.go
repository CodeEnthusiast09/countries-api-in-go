package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/CodeEnthusiast09/country-currency-api/internal/config"
	"github.com/CodeEnthusiast09/country-currency-api/internal/models"
)

type ExternalService struct {
	cfg *config.Config
}

func NewExternalService(cfg *config.Config) *ExternalService {
	return &ExternalService{cfg: cfg}
}

func (s *ExternalService) FetchCountries() ([]models.ExternalCountry, error) {
	resp, err := http.Get(s.cfg.CountriesAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach countries API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("countries API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read countries response: %w", err)
	}

	var countries []models.ExternalCountry
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, fmt.Errorf("failed to parse countries response: %w", err)
	}

	return countries, nil
}

func (s *ExternalService) FetchExchangeRates() (map[string]float64, error) {
	resp, err := http.Get(s.cfg.ExchangeRateAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach exchange rate API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exchange rate API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read exchange rate response: %w", err)
	}

	var result models.ExchangeRateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse exchange rate response: %w", err)
	}

	return result.Rates, nil
}
