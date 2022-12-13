package tree

import "fmt"

// FileDirEntry is a file or directory tree entry
type FileDirEntry interface {
	Name() string
	Size() int64
	Count() int
	Add(size int64)
}

// FileEntry is a file tree entry
type FileEntry struct {
	name  string
	size  int64
	count int
}

// NewFileEntry creates a new FileEntry
func NewFileEntry(name string, size int64) FileEntry {
	return FileEntry{
		name: name,
		size: size,
	}
}

// DirEntry is a directory tree entry
type DirEntry struct {
	FileEntry
	Extensions map[string]FileEntry
}

// NewDirEntry creates a new FileEntry
func NewDirEntry(name string) DirEntry {
	return DirEntry{
		FileEntry{
			name: name,
			size: 0,
		},
		map[string]FileEntry{},
	}
}

// Name returns the name of the entry
func (e FileEntry) Name() string { return e.name }

// Size returns the size of the entry
func (e FileEntry) Size() int64 { return e.size }

// Count returns the count of the entry
func (e FileEntry) Count() int { return e.count }

// Add adds size and a count of one
func (e *FileEntry) Add(size int64) {
	e.count++
	e.size += size
}

func (e FileEntry) String() string {
	return fmt.Sprintf(" %-16s %9d kB", e.Name(), e.Size()/1024)
}

func (e DirEntry) String() string {
	return fmt.Sprintf("-%-16s %9d kB (%d)", e.Name(), e.Size()/1024, e.Count())
}
