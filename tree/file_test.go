package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryCreate(t *testing.T) {
	file := NewFile("f", 1)
	dir := NewDir("d")

	assert.Equal(t, "f", file.Value.Name)
	assert.Equal(t, int64(1), file.Value.Size)
	assert.Equal(t, 1, file.Value.Count)

	assert.Equal(t, "d", dir.Value.Name)
	assert.Equal(t, int64(0), dir.Value.Size)
	assert.Equal(t, 0, dir.Value.Count)
}

func TestEntryAdd(t *testing.T) {
	dir := NewDir("d")

	dir.Value.Add(100)
	assert.Equal(t, int64(100), dir.Value.Size)
	assert.Equal(t, 1, dir.Value.Count)
}

func TestEntryAddExtensions(t *testing.T) {
	dir := NewDir("d")

	dir.Value.AddExtensions(map[string]*ExtensionEntry{".exe": {Name: ".exe", Size: 100, Count: 10}})
	assert.Equal(t, map[string]*ExtensionEntry{".exe": {Name: ".exe", Size: 100, Count: 10}}, dir.Value.Extensions)
}
