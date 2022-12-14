package util

import (
	"fmt"
	"math"
)

var memUnits []string = []string{
	"B",
	"kB",
	"MB",
	"GB",
	"TB",
	"PB",
	"EB",
	"ZB",
}

// FormatBytes formats bytes to B, kB, MB, ...
func FormatBytes(b int64) string {
	if b <= 0 {
		return "0 B"
	}

	exp := math.Log10(float64(b))
	exp1k := int(math.Round(exp*10000) / 10000 / 3)
	if exp1k == 0 {
		return fmt.Sprintf("%d B", b)
	}
	fac := math.Pow(10, float64(exp1k*3))
	value := float64(b) / fac
	if value >= 10 {
		return fmt.Sprintf("%d %s", int(value), memUnits[exp1k])
	}
	return fmt.Sprintf("%.1f %s", value, memUnits[exp1k])
}
