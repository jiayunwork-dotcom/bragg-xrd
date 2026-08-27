package cli

import (
	"fmt"
	"strconv"

	"bragg-xrd/internal/xrd"
)

func parseFloat(text, label string) (float64, error) {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number, got %q", label, text)
	}
	return value, nil
}

func parseHKL(text string) (xrd.HKL, error) {
	return xrd.MillersFromText(text)
}

func validateWavelength(value float64) error {
	return xrd.ValidateWavelength(value)
}

func validateLattice(value string) error {
	return xrd.ValidateLattice(value)
}

func requirePositive(value float64, label string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %g", label, value)
	}
	return nil
}

func normalize(text string) string {
	if len(text) > 0 && text[0] == '"' {
		if parsed, err := strconv.Unquote(text); err == nil {
			return parsed
		}
	}
	return text
}
