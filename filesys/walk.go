package filesys

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/gobwas/glob"
	"github.com/mlange-42/dirstat/tree"
)

// Walk searches through a directory tree
func Walk(dir string, exclude []string, maxDepth int, progres chan<- int64, done chan<- *tree.FileTree, erro chan<- error) {
	excludeGlobs := make([]glob.Glob, 0, len(exclude))
	for _, g := range exclude {
		excludeGlobs = append(excludeGlobs, glob.MustCompile(g))
	}

	anyFound := false

	t, err := walkDir(dir,
		func(path string, d fs.DirEntry, parent *tree.FileTree, depth int, err error) (*tree.FileTree, error) {
			if err != nil {
				erro <- err
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
			anyFound = true

			if !info.IsDir() {
				v := parent.Value
				ext := filepath.Ext(info.Name())
				if inf, ok := v.Extensions[ext]; ok {
					inf.AddMulti(info.Size(), 1)
				} else {
					e := tree.ExtensionEntry{Name: ext, Size: info.Size(), Count: 1}
					v.Extensions[ext] = &e
				}
				v.Add(info.Size())
			}

			progres <- info.Size()

			if maxDepth >= 0 && depth > maxDepth {
				return parent, nil
			}
			var subTree *tree.FileTree
			if info.IsDir() {
				subTree = tree.NewDir(info.Name())
			} else {
				subTree = tree.NewFile(info.Name(), info.Size())
			}

			if parent != nil {
				parent.AddTree(subTree)
			}
			return subTree, nil
		})

	if err != nil {
		erro <- err
		return
	}

	if !anyFound {
		erro <- fmt.Errorf("Nothing found in directoy %s", dir)
		return
	}

	t.Aggregate(func(parent, child *tree.FileEntry) {
		if child.IsDir {
			parent.AddMulti(child.Size, child.Count)
		}
	})

	done <- t
}

// walkDir recursively descends path, calling walkDirFn.
func walkDir[T any](root string, fn WalkDirFunc[T]) (*tree.Tree[T], error) {
	info, err := os.Lstat(root)
	var t *tree.Tree[T] = nil
	if err != nil {
		t, err = fn(root, nil, nil, 0, err)
	} else {
		t, err = walkDirRecursive(root, &statDirEntry{info}, nil, 0, fn)
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

// walkDirRecursive recursively descends path, calling walkDirFn.
func walkDirRecursive[T any](path string, d fs.DirEntry, parent *tree.Tree[T], depth int, walkDirFn WalkDirFunc[T]) (*tree.Tree[T], error) {
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
		_, err := walkDirRecursive(path1, d1, t, depth+1, walkDirFn)
		if err != nil && err != filepath.SkipDir {
			return nil, err
		}
	}
	return t, nil
}

// readDir reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDir(dirname string) ([]fs.DirEntry, error) {
	f, err := os.Open(dirname)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return []fs.DirEntry{}, nil
		}
		panic(err)
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(dirs,
		func(i, j int) bool {
			if dirs[i].IsDir() && !dirs[j].IsDir() {
				return true
			}
			if !dirs[i].IsDir() && dirs[j].IsDir() {
				return false
			}
			return lessCaseInsensitive(dirs[i].Name(), dirs[j].Name())
		})
	return dirs, nil
}

// lessCaseInsensitive compares s, t without allocating
func lessCaseInsensitive(s, t string) bool {
	for {
		if len(t) == 0 {
			return false
		}
		if len(s) == 0 {
			return true
		}
		c, sizec := utf8.DecodeRuneInString(s)
		d, sized := utf8.DecodeRuneInString(t)

		lowerc := unicode.ToLower(c)
		lowerd := unicode.ToLower(d)

		if lowerc < lowerd {
			return true
		}
		if lowerc > lowerd {
			return false
		}

		s = s[sizec:]
		t = t[sized:]
	}
}
