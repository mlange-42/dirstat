package util

import (
	"fmt"
	"math"
	"time"
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
		}
		return fmt.Sprintf("0 %s ", unit)
	}

	exp := math.Log10(float64(b))
	exp1k := int(math.Round(exp*10000) / 10000 / 3)
	if exp1k == 0 {
		if len(unit) == 0 {
			return fmt.Sprintf("%d  ", b)
		}
		return fmt.Sprintf("%d %s ", b, unit)
	}
	fac := math.Pow(10, float64(exp1k*3))
	value := float64(b) / fac
	if value >= 10 {
		return fmt.Sprintf("%d %s%s", int(value), unitPrefixes[exp1k], unit)
	}
	return fmt.Sprintf("%.1f %s%s", value, unitPrefixes[exp1k], unit)
}

// FormatUnitsSimple formats numbers with unit prefixes, like k, M, ..., without extra padding
func FormatUnitsSimple(b int64, unit string) string {
	if b <= 0 {
		if len(unit) == 0 {
			return fmt.Sprintf("0")
		}
		return fmt.Sprintf("0 %s", unit)
	}

	exp := math.Log10(float64(b))
	exp1k := int(math.Round(exp*10000) / 10000 / 3)
	if exp1k == 0 {
		if len(unit) == 0 {
			return fmt.Sprintf("%d", b)
		}
		return fmt.Sprintf("%d %s", b, unit)
	}
	fac := math.Pow(10, float64(exp1k*3))
	value := float64(b) / fac
	if value >= 10 {
		return fmt.Sprintf("%d %s%s", int(value), unitPrefixes[exp1k], unit)
	}
	return fmt.Sprintf("%.1f %s%s", value, unitPrefixes[exp1k], unit)
}

// FormatDuration prints a foratter duration to a Writer
func FormatDuration(from time.Time, to time.Time) string {
	if from.IsZero() || to.IsZero() {
		return "---"
	}
	dur := to.Sub(from)
	minutes := dur.Minutes()
	if minutes <= 60 {
		return fmt.Sprintf("%.0f minutes", minutes)
	}
	hours := dur.Hours()
	if hours <= 24 {
		return fmt.Sprintf("%.0f hours  ", hours)
	}
	days := hours / 24
	if days <= 14 {
		return fmt.Sprintf("%.0f days   ", days)
	}
	if days <= 60 {
		return fmt.Sprintf("%.0f weeks  ", days/7)
	}
	if days <= 2*365 {
		return fmt.Sprintf("%.0f months ", days/30.42)
	}
	years := days / 365
	return fmt.Sprintf("%.0f years  ", years)
}
