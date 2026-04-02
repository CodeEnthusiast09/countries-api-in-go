package main

import (
	"fmt"
	"log"

	"github.com/CodeEnthusiast09/country-currency-api/internal/config"
	"github.com/CodeEnthusiast09/country-currency-api/internal/database"
	"github.com/CodeEnthusiast09/country-currency-api/internal/router"
	"github.com/CodeEnthusiast09/country-currency-api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	db := database.New(cfg)

	externalService := services.NewExternalService(cfg)

	countryService := services.NewCountryService(db.DB, externalService)

	r := router.Setup(countryService)

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Starting server on port %s", cfg.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
