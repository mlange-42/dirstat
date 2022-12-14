package util

import (
	"fmt"
	"math"
)

var unitPrefixes []string = []string{
	"",
	"k",
	"M",
	"G",
	"T",
	"P",
	"E",
	"Z",
}

// FormatUnits formats numbers with unit prefixes, like k, M, ...
func FormatUnits(b int64, unit string) string {
	if b <= 0 {
		if len(unit) == 0 {
			return fmt.Sprintf("0")
		} else {
			return fmt.Sprintf("0 %s", unit)
		}
	}

	exp := math.Log10(float64(b))
	exp1k := int(math.Round(exp*10000) / 10000 / 3)
	if exp1k == 0 {
		if len(unit) == 0 {
			return fmt.Sprintf("%d", b)
		} else {
			return fmt.Sprintf("%d %s", b, unit)
		}
	}
	fac := math.Pow(10, float64(exp1k*3))
	value := float64(b) / fac
	if value >= 10 {
		return fmt.Sprintf("%d %s%s", int(value), unitPrefixes[exp1k], unit)
	}
	return fmt.Sprintf("%.1f %s%s", value, unitPrefixes[exp1k], unit)
}
