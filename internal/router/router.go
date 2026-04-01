package router

import (
	"github.com/CodeEnthusiast09/country-currency-api/internal/handlers"
	"github.com/CodeEnthusiast09/country-currency-api/internal/services"
	"github.com/gin-gonic/gin"
)

func Setup(countryService *services.CountryService) *gin.Engine {
	r := gin.Default()

	countryHandler := handlers.NewCountryHandler(countryService)

	countries := r.Group("/countries")
	{
		countries.POST("/refresh", countryHandler.Refresh)
		countries.GET("", countryHandler.GetAll)
		countries.GET("/image", countryHandler.GetImage)
		countries.GET("/:name", countryHandler.GetOne)
		countries.DELETE("/:name", countryHandler.Delete)
	}

	r.GET("/status", countryHandler.GetStatus)

	return r
}
