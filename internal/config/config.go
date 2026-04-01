package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	AtlasDatabaseURL   string
	DBName             string
	Port               string
	CountriesAPIURL    string
	ExchangeRateAPIURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment directly")
	}

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		AtlasDatabaseURL:   getEnv("ATLAS_DATABASE_URL", ""),
		DBName:             getEnv("DB_NAME", ""),
		Port:               getEnv("PORT", "3000"),
		CountriesAPIURL:    getEnv("COUNTRIES_API_URL", ""),
		ExchangeRateAPIURL: getEnv("EXCHANGE_RATE_API_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
