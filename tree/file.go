package tree

import (
	"fmt"

	"github.com/mlange-42/dirstat/util"
)

// FileTree is a tree with TreeEntry data
type FileTree = Tree[*FileEntry]

// NewDir creates a new FileTree with a directory entry
func NewDir(name string) *FileTree {
	e := NewFileEntry(name, 0, true)
	t := New(&e)
	return t
}

// NewFile creates a new FileTree with a file entry
func NewFile(name string, size int64) *FileTree {
	e := NewFileEntry(name, size, false)
	t := New(&e)
	return t
}

// GetSize returns the size
func (t *FileTree) GetSize() int64 {
	return t.Value.Size
}

// GetCount returns the count
func (t *FileTree) GetCount() int {
	return t.Value.Count
}

// Sized has a size
type Sized interface {
	GetSize() int64
}

// Counted has a count
type Counted interface {
	GetCount() int
}

// FileEntry is a file tree entry
type FileEntry struct {
	Name       string                     `json:"name"`
	IsDir      bool                       `json:"is_dir"`
	Size       int64                      `json:"size"`
	Count      int                        `json:"count"`
	Extensions map[string]*ExtensionEntry `json:"extensions"`
}

// ExtensionEntry is a file tree entry for extensions
type ExtensionEntry struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Count int    `json:"count"`
}

// NewFileEntry creates a new FileEntry
func NewFileEntry(name string, size int64, isDir bool) FileEntry {
	count := 0
	var ext map[string]*ExtensionEntry = nil
	if !isDir {
		count = 1
	} else {
		ext = map[string]*ExtensionEntry{}
	}
	return FileEntry{
		Name:       name,
		IsDir:      isDir,
		Size:       size,
		Count:      count,
		Extensions: ext,
	}
}

// GetSize returns the size
func (e *FileEntry) GetSize() int64 {
	return e.Size
}

// GetCount returns the count
func (e *FileEntry) GetCount() int {
	return e.Count
}

// Add adds size and a count of one
func (e *FileEntry) Add(size int64) {
	e.Count++
	e.Size += size
}

// AddMulti adds size and a count
func (e *FileEntry) AddMulti(size int64, count int) {
	e.Count += count
	e.Size += size
}

// GetSize returns the size
func (e *ExtensionEntry) GetSize() int64 {
	return e.Size
}

// GetCount returns the count
func (e *ExtensionEntry) GetCount() int {
	return e.Count
}

// AddMulti adds size and a count
func (e *ExtensionEntry) AddMulti(size int64, count int) {
	e.Count += count
	e.Size += size
}

func (e ExtensionEntry) String() string {
	return fmt.Sprintf("%s (%s)", util.FormatUnits(e.Size, "B"), util.FormatUnits(int64(e.Count), ""))
}

func (e FileEntry) String() string {
	if e.IsDir {
		return fmt.Sprintf("-%s %s (%s) %v", e.Name, util.FormatUnits(e.Size, "B"), util.FormatUnits(int64(e.Count), ""), e.Extensions)
	}
	return fmt.Sprintf(" %s %s (%s)", e.Name, util.FormatUnits(e.Size, "B"), util.FormatUnits(int64(e.Count), ""))
}

// AddExtensions adds extensions
func (e *FileEntry) AddExtensions(ext map[string]*ExtensionEntry) {
	for k, v := range ext {
		if inf, ok := e.Extensions[k]; ok {
			inf.AddMulti(v.Size, v.Count)
		} else {
			fe := ExtensionEntry{k, v.Size, v.Count}
			e.Extensions[k] = &fe
		}
	}
}
