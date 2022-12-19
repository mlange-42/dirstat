package tree

import (
	"fmt"
	"time"
	tm "time"

	"github.com/mlange-42/dirstat/util"
)

// FileTree is a tree with TreeEntry data
type FileTree = Tree[*FileEntry]

// NewDir creates a new FileTree with a directory entry
func NewDir(name string) *FileTree {
	e := NewFileEntry(name, 0, tm.Time{}, true)
	t := New(&e)
	return t
}

// NewFile creates a new FileTree with a file entry
func NewFile(name string, size int64, time tm.Time) *FileTree {
	e := NewFileEntry(name, size, time, false)
	t := New(&e)
	return t
}

// FileEntry is a file tree entry
type FileEntry struct {
	Name       string                     `json:"name"`
	IsDir      bool                       `json:"is_dir"`
	Size       int64                      `json:"size"`
	Count      int                        `json:"count"`
	Time       time.Time                  `json:"time"`
	Extensions map[string]*ExtensionEntry `json:"extensions"`
}

// ExtensionEntry is a file tree entry for extensions
type ExtensionEntry struct {
	Name  string    `json:"name"`
	Size  int64     `json:"size"`
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}

// NewFileEntry creates a new FileEntry
func NewFileEntry(name string, size int64, time tm.Time, isDir bool) FileEntry {
	count := 0
	var ext map[string]*ExtensionEntry = nil
	if isDir {
		ext = map[string]*ExtensionEntry{}
		time = tm.Time{}
	} else {
		count = 1
	}
	return FileEntry{
		Name:       name,
		IsDir:      isDir,
		Size:       size,
		Count:      count,
		Time:       time,
		Extensions: ext,
	}
}

// AddMulti adds size and a count
func (e *FileEntry) AddMulti(size int64, count int, time tm.Time) {
	e.Count += count
	e.Size += size
	if !time.IsZero() && (e.Time.IsZero() || time.After(e.Time)) {
		e.Time = time
	}
}

// AddMulti adds size and a count
func (e *ExtensionEntry) AddMulti(size int64, count int, time tm.Time) {
	e.Count += count
	e.Size += size
	if !time.IsZero() && (e.Time.IsZero() || time.After(e.Time)) {
		e.Time = time
	}
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
			inf.AddMulti(v.Size, v.Count, v.Time)
		} else {
			fe := ExtensionEntry{k, v.Size, v.Count, v.Time}
			e.Extensions[k] = &fe
		}
	}
}
