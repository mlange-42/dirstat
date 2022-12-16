package tree

import (
	"fmt"
)

// Tree is a tree data structure
type Tree[T any] struct {
	Children []*Tree[T] `json:"children"`
	Value    T          `json:"value"`
}

// New creates a new tree
func New[T any](value T) *Tree[T] {
	return &Tree[T]{
		Children: make([]*Tree[T], 0, 0),
		Value:    value,
	}
}

// AddTree adds a sub-tree without children
func (t *Tree[T]) Add(child T) {
	t.Children = append(t.Children, New(child))
}

// AddTree adds a sub-tree
func (t *Tree[T]) AddTree(child *Tree[T]) {
	t.Children = append(t.Children, child)
}

// Aggregate aggregates tree branches
func (t *Tree[T]) Aggregate(fn func(parent, child T)) {
	for _, child := range t.Children {
		child.Aggregate(fn)
		fn(t.Value, child.Value)
	}
}

// Crop reduces the tree to a given depth
func (t *Tree[T]) Crop(depth int, fn func(parent, child T)) {
	if depth <= 0 {
		t.Aggregate(fn)
		t.Children = nil
		return
	}
	for _, child := range t.Children {
		child.Crop(depth-1, fn)
	}
}

// String converts the tree to a multiline string
func (t *Tree[T]) String() string {
	return PlainPrinter[T]{}.Print(t)
}

// SubTree selects a sub-tree using a path
func SubTree[T any, P any](t *Tree[T], path []P, fn func(T, P) bool) (*Tree[T], error) {
	if len(path) == 0 {
		return t, nil
	}
	for _, child := range t.Children {
		if fn(child.Value, path[0]) {
			if len(path) == 1 {
				return child, nil
			}
			return SubTree(child, path[1:], fn)
		}
	}
	return nil, fmt.Errorf("Path element '%v' not found in tree", path[0])
}
