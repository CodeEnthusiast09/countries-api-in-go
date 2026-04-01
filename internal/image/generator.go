package image

import (
	"fmt"
	"time"

	"github.com/CodeEnthusiast09/country-currency-api/internal/lib"
	"github.com/CodeEnthusiast09/country-currency-api/internal/models"
	"github.com/fogleman/gg"
)

const (
	imgWidth  = 800
	imgHeight = 400
	fontPath  = "/usr/share/fonts/TTF/DejaVuSans.ttf"
)

func Generate(total int64, topCountries []models.Country, lastRefreshed time.Time) (string, error) {
	dc := gg.NewContext(imgWidth, imgHeight)

	// background
	dc.SetHexColor("#1a1a2e")
	dc.Clear()

	// title
	if err := dc.LoadFontFace(fontPath, 28); err != nil {
		return "", fmt.Errorf("failed to load font: %w", err)
	}
	dc.SetHexColor("#e0e0e0")
	dc.DrawStringAnchored("Country Currency & Exchange Summary", imgWidth/2, 40, 0.5, 0.5)

	// total countries
	dc.LoadFontFace(fontPath, 18)
	dc.SetHexColor("#a0a0ff")
	dc.DrawString(fmt.Sprintf("Total Countries: %d", total), 40, 100)

	// last refreshed
	dc.LoadFontFace(fontPath, 14)
	dc.SetHexColor("#808080")
	dc.DrawString(fmt.Sprintf("Last Refreshed: %s", lastRefreshed.Format("2006-01-02 15:04:05")), 40, 130)

	// top 5 header
	dc.LoadFontFace(fontPath, 18)
	dc.SetHexColor("#e0e0e0")
	dc.DrawString("Top 5 Countries by Estimated GDP:", 40, 180)

	// top 5 list
	dc.LoadFontFace(fontPath, 15)
	for i, c := range topCountries {
		gdp := 0.0

		if c.EstimatedGDP != nil {
			gdp = *c.EstimatedGDP
		}

		formattedGDP := lib.FormatNumberWithSuffix(gdp)

		line := fmt.Sprintf("%d. %s — $%s", i+1, c.Name, formattedGDP)

		dc.SetHexColor("#c0c0c0")

		dc.DrawString(line, 60, float64(210+(i*35)))
	}

	outputPath := "/tmp/summary.png"
	if err := dc.SavePNG(outputPath); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return outputPath, nil
}
