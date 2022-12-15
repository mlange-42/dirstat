package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeCreate(t *testing.T) {
	tr := New(1)

	assert.Equal(t, 1, tr.Value)
	assert.Equal(t, 0, len(tr.Children))
}

func TestTreeAdd(t *testing.T) {
	tr := New(1)

	tr.Add(2)
	tr.Add(3)

	assert.Equal(t, 1, tr.Value)
	assert.Equal(t, 2, len(tr.Children))
	assert.Equal(t, 2, tr.Children[0].Value)
	assert.Equal(t, 3, tr.Children[1].Value)
}

func TestTreeAddTree(t *testing.T) {
	tr := New(1)

	tr.AddTree(New(2))
	tr.AddTree(New(3))

	assert.Equal(t, 1, tr.Value)
	assert.Equal(t, 2, len(tr.Children))
	assert.Equal(t, 2, tr.Children[0].Value)
	assert.Equal(t, 3, tr.Children[1].Value)
}

func TestTreeAggregate(t *testing.T) {
	a := 1
	b := 2
	c := 3

	tr := New(&a)

	tr.AddTree(New(&b))
	tr.AddTree(New(&c))

	tr.Aggregate(func(parent, child *int) {
		*parent += *child
	})

	assert.Equal(t, 6, *tr.Value)
}

func TestTreeCrop(t *testing.T) {
	a := 1
	b := 2
	c := 3

	tr := New(&a)

	tr.AddTree(New(&b))
	tr.AddTree(New(&c))

	tr.Crop(0, func(parent, child *int) {
		*parent += *child
	})

	assert.Equal(t, 6, *tr.Value)
	assert.Equal(t, 0, len(tr.Children))
}
