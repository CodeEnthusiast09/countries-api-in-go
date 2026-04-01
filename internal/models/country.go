package models

import "time"

type Country struct {
	ID              uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	Name            string    `json:"name"             gorm:"column:name;not null;uniqueIndex"`
	Capital         string    `json:"capital"          gorm:"column:capital"`
	Region          string    `json:"region"           gorm:"column:region"`
	Population      int64     `json:"population"       gorm:"column:population"`
	CurrencyCode    *string   `json:"currency_code"     gorm:"column:currency_code"`
	ExchangeRate    *float64  `json:"exchange_rate"     gorm:"column:exchange_rate"`
	EstimatedGDP    *float64  `json:"estimated_gdp"     gorm:"column:estimated_gdp"`
	FlagURL         string    `json:"flag_url"          gorm:"column:flag_url"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"  gorm:"column:last_refreshed_at"`
}
