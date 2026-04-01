package handlers

import (
	"net/http"
	"time"

	"github.com/CodeEnthusiast09/country-currency-api/internal/image"
	"github.com/CodeEnthusiast09/country-currency-api/internal/services"
	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	countryService *services.CountryService
}

func NewCountryHandler(countryService *services.CountryService) *CountryHandler {
	return &CountryHandler{countryService: countryService}
}

func (h *CountryHandler) Refresh(c *gin.Context) {
	total, err := h.countryService.Refresh()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		Message: "Countries data refreshed successfully",
		Total:   total,
	})
}

func (h *CountryHandler) GetAll(c *gin.Context) {
	region := c.Query("region")
	currency := c.Query("currency")
	sort := c.Query("sort")

	countries, err := h.countryService.GetAll(region, currency, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, countries)
}

func (h *CountryHandler) GetImage(c *gin.Context) {
	status, err := h.countryService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	topCountries, err := h.countryService.GetTopCountries(5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	imgPath, err := image.Generate(status.Total, topCountries, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate image"})
		return
	}

	c.File(imgPath)
}

func (h *CountryHandler) GetOne(c *gin.Context) {
	name := c.Param("name")

	country, err := h.countryService.GetOne(name)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Country not found"})
		return
	}

	c.JSON(http.StatusOK, country)
}

func (h *CountryHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	if err := h.countryService.Delete(name); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Country not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Country deleted successfully"})
}

func (h *CountryHandler) GetStatus(c *gin.Context) {
	status, err := h.countryService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, status)
}
