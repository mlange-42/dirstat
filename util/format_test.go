package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	assert.Equal(t, "0 B", FormatBytes(0))
	assert.Equal(t, "123 B", FormatBytes(123))
	assert.Equal(t, "1.2 kB", FormatBytes(1234))
	assert.Equal(t, "12 kB", FormatBytes(12345))
	assert.Equal(t, "1.2 MB", FormatBytes(1234567))
	assert.Equal(t, "1.0 MB", FormatBytes(1e6))
	assert.Equal(t, "1.0 GB", FormatBytes(1e9))
	assert.Equal(t, "1.0 TB", FormatBytes(1e12))
	assert.Equal(t, "1.0 PB", FormatBytes(1e15))
}
