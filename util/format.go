package util

import (
	"fmt"
	"io"
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

// FPrintDuration prints a foratter duration to a Writer
func FPrintDuration(w io.Writer, dur time.Duration) {
	minutes := dur.Minutes()
	if minutes <= 60 {
		fmt.Fprintf(w, "%.0f minutes", minutes)
		return
	}
	hours := dur.Hours()
	if hours <= 24 {
		fmt.Fprintf(w, "%.0f hours", hours)
		return
	}
	days := hours / 24
	if days <= 14 {
		fmt.Fprintf(w, "%.0f days", days)
		return
	}
	if days <= 60 {
		fmt.Fprintf(w, "%.0f weeks", days/7)
		return
	}
	if days <= 24 {
		fmt.Fprintf(w, "%.0f months", days/30.42)
		return
	}
	years := days / 365
	fmt.Fprintf(w, "%.0f years", years)
}
