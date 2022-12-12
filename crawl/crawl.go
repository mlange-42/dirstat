package crawl

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/gobwas/glob"
	"github.com/mlange42/dirstat/tree"
)

// FileTree is a tree with TreeEntry data
type FileTree = tree.Tree[TreeEntry]

// NewFileTree creates a new FileTree
func NewFileTree(name string, size int64, isDir bool) *FileTree {
	t := tree.New(TreeEntry{
		Name:  name,
		Size:  size,
		Count: 0,
		IsDir: isDir,
	})
	return t
}

// TreeEntry is a file tree entry
type TreeEntry struct {
	Name  string
	Count int
	Size  int64
	IsDir bool
}

func (e TreeEntry) String() string {
	pref := " "
	if e.IsDir {
		pref = "-"
	}
	return fmt.Sprintf("%s%-16s %9d kB (%d)", pref, e.Name, e.Size/1024, e.Count)
}

// Walk searches through a directory tree
func Walk(dir string, exclude []string, maxDepth int) (*tree.Tree[TreeEntry], error) {
	excludeGlobs := make([]glob.Glob, 0, len(exclude))
	for _, g := range exclude {
		excludeGlobs = append(excludeGlobs, glob.MustCompile(g))
	}

	dir = path.Clean(dir)
	//fsMap := map[string]fs.FileInfo{}

	t, err := WalkDir(dir,
		func(path string, d fs.DirEntry, parent *FileTree, depth int, err error) (*FileTree, error) {
			if err != nil {
				return nil, err
			}
			for _, g := range excludeGlobs {
				if g.Match(d.Name()) {
					return nil, fs.SkipDir
				}
			}
			info, err := d.Info()
			if err != nil {
				return nil, err
			}

			if maxDepth >= 0 && depth > maxDepth {
				parent.Value.Count++
				parent.Value.Size += info.Size()
				return parent, nil
			}

			subTree := NewFileTree(info.Name(), info.Size(), info.IsDir())
			if parent != nil {
				parent.AddTree(subTree)
			}
			return subTree, nil
		})
	if err != nil {
		return nil, err
	}
	return t, nil
}

// WalkDir recursively descends path, calling walkDirFn.
func WalkDir[T any](root string, fn WalkDirFunc[T]) (*tree.Tree[T], error) {
	info, err := os.Lstat(root)
	var t *tree.Tree[T] = nil
	if err != nil {
		t, err = fn(root, nil, nil, 0, err)
	} else {
		t, err = walkDir(root, &statDirEntry{info}, nil, 0, fn)
	}
	if err == filepath.SkipDir {
		return t, nil
	}
	return t, err
}

// WalkDirFunc as callback for WalkDir
type WalkDirFunc[T any] func(path string, d fs.DirEntry, parent *tree.Tree[T], depth int, err error) (*tree.Tree[T], error)

type statDirEntry struct {
	info fs.FileInfo
}

func (d *statDirEntry) Name() string               { return d.info.Name() }
func (d *statDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *statDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *statDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

// walkDir recursively descends path, calling walkDirFn.
func walkDir[T any](path string, d fs.DirEntry, parent *tree.Tree[T], depth int, walkDirFn WalkDirFunc[T]) (*tree.Tree[T], error) {
	t, err := walkDirFn(path, d, parent, depth, nil)
	if err != nil || !d.IsDir() {
		if err == filepath.SkipDir {
			// Successfully skipped directory.
			err = nil
		}
		return t, err
	}

	dirs, err := readDir(path)
	if err != nil {
		// Second call, to report ReadDir error.
		_, err = walkDirFn(path, d, parent, depth, err)
		if err != nil {
			if err == filepath.SkipDir && d.IsDir() {
				err = nil
			}
			return t, err
		}
	}

	for _, d1 := range dirs {
		path1 := filepath.Join(path, d1.Name())
		_, err := walkDir(path1, d1, t, depth+1, walkDirFn)
		if err != nil {
			if err == filepath.SkipDir {
				break
			}
			return t, err
		}
	}
	return t, nil
}

// readDir reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDir(dirname string) ([]fs.DirEntry, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	return dirs, nil
}
