package models

import "time"

type Country struct {
	ID              uint      `json:"id"               gorm:"primaryKey;autoIncrement"`
	Name            string    `json:"name"              gorm:"column:name;size:255;not null;uniqueIndex"`
	Capital         string    `json:"capital"           gorm:"column:capital;size:255;not null;default:''"`
	Region          string    `json:"region"            gorm:"column:region;size:100;not null;default:''"`
	Population      int64     `json:"population"        gorm:"column:population;not null;default:0"`
	CurrencyCode    *string   `json:"currency_code"     gorm:"column:currency_code;size:10"`
	ExchangeRate    *float64  `json:"exchange_rate"     gorm:"column:exchange_rate"`
	EstimatedGDP    *float64  `json:"estimated_gdp"     gorm:"column:estimated_gdp"`
	FlagURL         string    `json:"flag_url"          gorm:"column:flag_url;size:500;not null;default:''"`
	LastRefreshedAt time.Time `json:"last_refreshed_at" gorm:"column:last_refreshed_at;not null"`
}
