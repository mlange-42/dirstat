package tree

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Printer interface
type Printer[T any] interface {
	Print(t *Tree[T]) string
}

// JSONPrinter prints a tree in JSON format
type JSONPrinter[T any] struct{}

func (p JSONPrinter[T]) Print(t *Tree[T]) string {
	tt, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(tt[:])
}

// PlainPrinter prints a tree in plain text format
type PlainPrinter[T any] struct{}

func (p PlainPrinter[T]) Print(t *Tree[T]) string {
	sb := strings.Builder{}
	p.print(t, &sb, 0)
	return sb.String()
}

func (p PlainPrinter[T]) print(t *Tree[T], sb *strings.Builder, depth int) {
	fmt.Fprint(sb, strings.Repeat(" ", depth*2))
	fmt.Fprintf(sb, "%v\n", t.Value)
	for _, child := range t.Children {
		p.print(child, sb, depth+1)
	}
}

// TreemapPrinter prints a tree in treemap CSV format
type TreemapPrinter struct {
	ByExtension bool
}

// Print prints a FileTree
func (p TreemapPrinter) Print(t *FileTree) string {
	sb := strings.Builder{}
	p.print(t, &sb, "")
	return sb.String()
}

func (p TreemapPrinter) print(t *FileTree, sb *strings.Builder, path string) {
	if len(path) == 0 {
		path = t.Value.Name
	} else {
		path += "/" + t.Value.Name
	}
	fmt.Fprintf(sb, "%s,%d,%d\n", strings.Replace(path, ",", "-", -1), t.Value.Size, t.Value.Count)
	if p.ByExtension && t.Value.IsDir && len(t.Children) == 0 {
		for _, info := range t.Value.Extensions {
			p := path + "/" + info.Name
			fmt.Fprintf(sb, "%s,%d,%d\n", p, info.Size, info.Count)
		}
		return
	}
	for _, child := range t.Children {
		p.print(child, sb, path)
	}
}
