package main

import (
	"testing"

	"github.com/mlange42/dirstat/crawl"
)

func TestRun(*testing.T) {
	t, _ := crawl.Walk("..", []string{".git"}, -1)

	_ = t.Value
}
