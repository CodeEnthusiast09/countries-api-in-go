package main

import (
	"fmt"
	"log"

	"github.com/CodeEnthusiast09/country-currency-api/internal/config"
	"github.com/CodeEnthusiast09/country-currency-api/internal/database"
	"github.com/CodeEnthusiast09/country-currency-api/internal/router"
	"github.com/CodeEnthusiast09/country-currency-api/internal/services"
)

func main() {
	cfg := config.Load()

	db := database.New(cfg)

	externalService := services.NewExternalService(cfg)

	countryService := services.NewCountryService(db.DB, externalService)

	r := router.Setup(countryService)

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Starting server on port %s", cfg.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
