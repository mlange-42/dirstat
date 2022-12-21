package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatUnitsSimple(t *testing.T) {
	assert.Equal(t, "0 B", FormatUnitsSimple(0, "B"))
	assert.Equal(t, "123 B", FormatUnitsSimple(123, "B"))
	assert.Equal(t, "1.2 kB", FormatUnitsSimple(1234, "B"))
	assert.Equal(t, "12 kB", FormatUnitsSimple(12345, "B"))
	assert.Equal(t, "1.2 MB", FormatUnitsSimple(1234567, "B"))
	assert.Equal(t, "1.0 MB", FormatUnitsSimple(1e6, "B"))
	assert.Equal(t, "1.0 GB", FormatUnitsSimple(1e9, "B"))
	assert.Equal(t, "1.0 TB", FormatUnitsSimple(1e12, "B"))
	assert.Equal(t, "1.0 PB", FormatUnitsSimple(1e15, "B"))
}

func TestFormatUnits(t *testing.T) {
	assert.Equal(t, "0 B ", FormatUnits(0, "B"))
	assert.Equal(t, "123 B ", FormatUnits(123, "B"))
	assert.Equal(t, "1.2 kB", FormatUnits(1234, "B"))
	assert.Equal(t, "12 kB", FormatUnits(12345, "B"))
	assert.Equal(t, "1.2 MB", FormatUnits(1234567, "B"))
	assert.Equal(t, "1.0 MB", FormatUnits(1e6, "B"))
	assert.Equal(t, "1.0 GB", FormatUnits(1e9, "B"))
	assert.Equal(t, "1.0 TB", FormatUnits(1e12, "B"))
	assert.Equal(t, "1.0 PB", FormatUnits(1e15, "B"))
}
