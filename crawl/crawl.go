package crawl

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/mlange42/dirstat/tree"
)

// FileTree is a tree with TreeEntry data
type FileTree = tree.Tree[TreeEntry]

// NewFileTree creates a new FileTree
func NewFileTree() *FileTree {
	t := tree.New(TreeEntry{})
	return t
}

// TreeEntry is a file tree entry
type TreeEntry struct {
	Name  string
	Count int
	Size  int
}

// Walk searches through a directory tree
func Walk(dir string, exclude map[string]bool) (*tree.Tree[TreeEntry], error) {
	dir = path.Clean(dir)
	//fsMap := map[string]fs.FileInfo{}

	t, err := WalkDir(dir,
		func(path string, d fs.DirEntry, parent *FileTree, depth int, err error) (*FileTree, error) {
			if err != nil {
				return nil, err
			}
			if _, ok := exclude[d.Name()]; ok {
				return nil, fs.SkipDir
			}
			_, err = d.Info()
			if err != nil {
				return nil, err
			}
			subTree := NewFileTree()
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
		if err == filepath.SkipDir && d.IsDir() {
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
