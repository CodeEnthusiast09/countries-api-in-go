package lib

import (
	"fmt"
	"math"

	"golang.org/x/text/message"
)

var printer = message.NewPrinter(message.MatchLanguage("en"))

func FormatNumberWithSuffix(num float64) string {
	if math.IsNaN(num) {
		return "Invalid amount"
	}

	thousand := 1_000.0
	million := 1_000_000.0
	billion := 1_000_000_000.0
	trillion := 1_000_000_000_000.0

	switch {
	case num < thousand:
		return formatWithCommas(num)
	case num < million:
		return formatLargeNumber(num, thousand, "K")
	case num < billion:
		return formatLargeNumber(num, million, "M")
	case num < trillion:
		return formatLargeNumber(num, billion, "B")
	default:
		return formatLargeNumber(num, trillion, "T")
	}
}

func formatLargeNumber(value float64, divisor float64, suffix string) string {
	dividedValue := float64(value) / float64(divisor)

	if math.Mod(dividedValue, 1) == 0 {
		return fmt.Sprintf("%.0f%s", dividedValue, suffix)
	} else {
		return fmt.Sprintf("%.1f%s", dividedValue, suffix)
	}
}

func formatWithCommas(num float64) string {
	return printer.Sprintf("%.0f", num)
}
