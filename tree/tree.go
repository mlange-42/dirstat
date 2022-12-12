package tree

// Tree is a tree data structure
type Tree[T any] struct {
	Children []*Tree[T]
	Value    T
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
