package tree

import "fmt"

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

// AddMulti adds size and a count
func (e *ExtensionEntry) AddMulti(size int64, count int) {
	e.Count += count
	e.Size += size
}

func (e ExtensionEntry) String() string {
	return fmt.Sprintf("%d kB (%d)", e.Size/1024, e.Count)
}

func (e FileEntry) String() string {
	if e.IsDir {
		return fmt.Sprintf("-%s %d kB (%d) %v", e.Name, e.Size/1024, e.Count, e.Extensions)
	} else {
		return fmt.Sprintf(" %s %d kB (%d)", e.Name, e.Size/1024, e.Count)
	}
}

// AddExtensions adds extensions
func (e *FileEntry) AddExtensions(ext map[string]*ExtensionEntry) {
	for k, v := range ext {
		if inf, ok := e.Extensions[k]; ok {
			inf.AddMulti(v.Size, v.Count)
		} else {
			fe := ExtensionEntry{k, v.Size, 1}
			e.Extensions[k] = &fe
		}
	}
}
