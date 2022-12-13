package tree

import "fmt"

// FileDirEntry is a file or directory tree entry
type FileDirEntry interface {
	GetName() string
	GetSize() int64
	GetCount() int
	Add(size int64)
}

// FileEntry is a file tree entry
type FileEntry struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Count int    `json:"count"`
}

// NewFileEntry creates a new FileEntry
func NewFileEntry(name string, size int64) FileEntry {
	return FileEntry{
		Name:  name,
		Size:  size,
		Count: 1,
	}
}

// DirEntry is a directory tree entry
type DirEntry struct {
	FileEntry
	Extensions map[string]*FileEntry `json:"extensions"`
}

// NewDirEntry creates a new FileEntry
func NewDirEntry(name string) DirEntry {
	return DirEntry{
		FileEntry{
			Name:  name,
			Size:  0,
			Count: 0,
		},
		map[string]*FileEntry{},
	}
}

// GetName returns the name of the entry
func (e FileEntry) GetName() string { return e.Name }

// GetSize returns the size of the entry
func (e FileEntry) GetSize() int64 { return e.Size }

// GetCount returns the count of the entry
func (e FileEntry) GetCount() int { return e.Count }

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

func (e FileEntry) String() string {
	return fmt.Sprintf(" %s %d kB (%d)", e.GetName(), e.GetSize()/1024, e.GetCount())
}

func (e DirEntry) String() string {
	return fmt.Sprintf("-%s %d kB (%d) %v", e.GetName(), e.GetSize()/1024, e.GetCount(), e.Extensions)
}

// AddExtensions adds extensions
func (e *DirEntry) AddExtensions(ext map[string]*FileEntry) {
	for k, v := range ext {
		if inf, ok := e.Extensions[k]; ok {
			inf.AddMulti(v.GetSize(), v.GetCount())
		} else {
			fe := NewFileEntry(k, v.GetSize())
			e.Extensions[k] = &fe
		}
	}
}
