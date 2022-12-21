package tree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntryCreate(t *testing.T) {
	file := NewFile("f", 1, time.Time{})
	dir := NewDir("d")

	assert.Equal(t, "f", file.Value.Name)
	assert.Equal(t, int64(1), file.Value.Size)
	assert.Equal(t, 1, file.Value.Count)

	assert.Equal(t, "d", dir.Value.Name)
	assert.Equal(t, int64(0), dir.Value.Size)
	assert.Equal(t, 0, dir.Value.Count)
}

func TestEntryAdd(t *testing.T) {
	tm := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)

	dir := NewDir("d")
	assert.Equal(t, true, dir.Value.Time.IsZero())

	dir.Value.Add(100, 1, tm)
	assert.Equal(t, int64(100), dir.Value.Size)
	assert.Equal(t, 1, dir.Value.Count)
	assert.Equal(t, tm, dir.Value.Time)
}

func TestEntryAddExtensions(t *testing.T) {
	tm := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)

	dir := NewDir("d")

	dir.Value.AddExtensions(map[string]*ExtensionEntry{".exe": {Name: ".exe", Size: 100, Count: 10, Time: tm}})
	assert.Equal(t, map[string]*ExtensionEntry{".exe": {Name: ".exe", Size: 100, Count: 10, Time: tm}}, dir.Value.Extensions)
}
