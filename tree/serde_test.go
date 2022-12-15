package tree

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {
	tr := NewDir("root")
	tr.AddTree(NewDir("b"))
	tr.AddTree(NewFile("c", 100))
	tr.Children[0].AddTree(NewDir("d"))

	b, err := json.MarshalIndent(tr, "", "    ")
	assert.Equal(t, nil, err)

	tr2, err := Deserialize(b)
	assert.Equal(t, nil, err)

	assert.Equal(t, tr, tr2)
}
