package tree

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mlange-42/go-dirstat/util"
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
	ByCount     bool
}

// Print prints a FileTree
func (p TreemapPrinter) Print(t *FileTree) string {
	sb := strings.Builder{}
	p.print(t, &sb, "")
	return sb.String()
}

func (p TreemapPrinter) print(t *FileTree, sb *strings.Builder, path string) {
	if len(path) == 0 {
		path = fmt.Sprintf(
			"%s (%s | %s)", t.Value.Name,
			util.FormatUnits(t.Value.Size, "B"),
			util.FormatUnits(int64(t.Value.Count), ""),
		)
	} else {
		path = fmt.Sprintf(
			"%s/%s (%s | %s)", path, t.Value.Name,
			util.FormatUnits(t.Value.Size, "B"),
			util.FormatUnits(int64(t.Value.Count), ""),
		)
	}

	var v1, v2 int64
	if p.ByCount {
		v1, v2 = int64(t.Value.Count), t.Value.Size
	} else {
		v1, v2 = t.Value.Size, int64(t.Value.Count)
	}

	fmt.Fprintf(
		sb,
		"%s,%d,%d\n",
		strings.Replace(path, ",", "-", -1),
		v1,
		v2,
	)

	if p.ByExtension && t.Value.IsDir && len(t.Children) == 0 {
		for _, info := range t.Value.Extensions {
			pth := path + "/" + info.Name
			if p.ByCount {
				v1, v2 = int64(info.Count), info.Size
			} else {
				v1, v2 = info.Size, int64(info.Count)
			}
			fmt.Fprintf(
				sb,
				"%s (%s | %s),%d,%d\n",
				strings.Replace(pth, ",", "-", -1),
				util.FormatUnits(info.Size, "B"),
				util.FormatUnits(int64(info.Count), ""),
				v1,
				v2,
			)
		}
		return
	}
	for _, child := range t.Children {
		p.print(child, sb, path)
	}
}
