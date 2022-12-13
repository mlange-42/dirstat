package tree

import (
	"fmt"
	"strings"
)

// Tree is a tree data structure
type Tree[T any] struct {
	Children []*Tree[T] `json:"children"`
	Value    T          `json:"value"`
}

// New creates a new tree
func New[T any](value T) *Tree[T] {
	return &Tree[T]{
		Value: value,
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

// String converts the tree to a multiline string
func (t *Tree[T]) String() string {
	sb := strings.Builder{}
	t.toString(&sb, 0)
	return sb.String()
}

func (t *Tree[T]) toString(sb *strings.Builder, depth int) {
	fmt.Fprint(sb, strings.Repeat(" ", depth*2))
	fmt.Fprintf(sb, "%v\n", t.Value)
	for _, child := range t.Children {
		child.toString(sb, depth+1)
	}
}
