package main

import (
	"testing"

	"github.com/mlange42/dirstat/filesys"
)

func TestRun(*testing.T) {
	t, _ := filesys.Walk("..", []string{".git"}, -1)

	_ = t.Value
}
