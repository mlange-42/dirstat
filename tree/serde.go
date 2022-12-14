package tree

import "encoding/json"

// Deserialize a FileTree from a JSON byte slice
func Deserialize(bytes []byte) (*FileTree, error) {
	t := NewDir("root")
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
