package tree

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/mlange-42/dirstat/util"
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

// FileTreePrinter prints a file tree in plain text format
type FileTreePrinter struct {
	ByExtension  bool
	Indent       int
	prefixNone   string
	prefixEmpty  string
	prefixNormal string
	prefixLast   string
	printWidth   int
}

// NewFileTreePrinter creates a new FileTreePrinter
func NewFileTreePrinter(byExt bool, indent int) FileTreePrinter {
	return FileTreePrinter{
		ByExtension:  byExt,
		Indent:       indent,
		prefixNone:   strings.Repeat(" ", indent),
		prefixEmpty:  "│" + strings.Repeat(" ", indent-1),
		prefixNormal: "├" + strings.Repeat("─", indent-1),
		prefixLast:   "└" + strings.Repeat("─", indent-1),
		printWidth:   0,
	}
}

// Print prints a FileTree
func (p FileTreePrinter) Print(t *FileTree) string {
	p.printWidth = p.maxWidth(t, 0) + 2
	if p.printWidth < 16 {
		p.printWidth = 16
	} else if p.printWidth > 64 {
		p.printWidth = 64
	}

	sb := strings.Builder{}
	p.print(t, &sb, 0, false, "")
	return sb.String()
}

func (p FileTreePrinter) print(t *FileTree, sb *strings.Builder, depth int, last bool, prefix string) {
	var sizeCount string
	if t.Value.IsDir {
		sizeCount = fmt.Sprintf("%-7s (%s)",
			util.FormatUnits(t.Value.Size, "B"), util.FormatUnits(int64(t.Value.Count), ""),
		)
	} else {
		sizeCount = fmt.Sprintf("%s", util.FormatUnits(t.Value.Size, "B"))
	}

	pref := prefix

	if depth > 0 {
		pref = prefix + p.createPrefixSimple(last)
	}
	pad := strings.Repeat(".", int(math.Max(float64(p.printWidth-depth*p.Indent-len([]rune(t.Value.Name))), 0)))
	fmt.Fprint(sb, pref)
	if t.Value.IsDir {
		fmt.Fprintf(sb, "%s/%s%s\n", t.Value.Name, pad, sizeCount)
	} else {
		fmt.Fprintf(sb, "%s %s%s\n", t.Value.Name, pad, sizeCount)
	}

	if depth > 0 {
		pref = prefix + p.createPrefixEmpty(last)
	}
	for i, child := range t.Children {
		if !p.ByExtension || child.Value.IsDir {
			p.print(child, sb, depth+1, i == len(t.Children)-1, pref)
		}
	}
	if p.ByExtension && t.Value.IsDir {
		p.printExtensions(t.Value.Extensions, sb, depth+1, pref)
	}
}

func (p FileTreePrinter) printExtensions(ext map[string]*ExtensionEntry, sb *strings.Builder, depth int, prefix string) {
	keys := make([]string, 0, len(ext))
	for k := range ext {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	prefix = prefix + p.createPrefixSimple(false)
	prefixLast := p.createPrefixSimple(true)

	for i, name := range keys {
		info := ext[name]

		pad := strings.Repeat(" ", int(math.Max(float64(p.printWidth-(depth)*p.Indent-len([]rune(info.Name))), 0)))
		if i == len(keys)-1 {
			fmt.Fprint(sb, prefixLast)
		} else {
			fmt.Fprint(sb, prefix)
		}
		fmt.Fprintf(
			sb,
			"%s %s%-7s (%s)\n",
			info.Name,
			pad,
			util.FormatUnits(info.Size, "B"),
			util.FormatUnits(int64(info.Count), ""),
		)
	}
}

func (p FileTreePrinter) maxWidth(t *FileTree, depth int) int {
	max := len([]rune(t.Value.Name)) + depth*p.Indent
	for _, c := range t.Children {
		v := p.maxWidth(c, depth+1)
		if v > max {
			max = v
		}
	}
	return max
}

func (p FileTreePrinter) createPrefix(depth, indent int, last bool) string {
	if depth <= 0 {
		return ""
	}
	if last {
		return strings.Repeat(p.prefixEmpty, depth-1) + p.prefixLast
	}
	return strings.Repeat(p.prefixEmpty, depth-1) + p.prefixNormal
}

func (p FileTreePrinter) createPrefixSimple(last bool) string {
	if last {
		return p.prefixLast
	}
	return p.prefixNormal
}

func (p FileTreePrinter) createPrefixEmpty(last bool) string {
	if last {
		return p.prefixNone
	}
	return p.prefixEmpty
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
	var sizeCount string

	if t.Value.IsDir {
		sizeCount = fmt.Sprintf("%s | %s",
			util.FormatUnits(t.Value.Size, "B"), util.FormatUnits(int64(t.Value.Count), ""),
		)
	} else {
		sizeCount = fmt.Sprintf("%s", util.FormatUnits(t.Value.Size, "B"))
	}

	if len(path) == 0 {
		path = fmt.Sprintf("%s (%s)", t.Value.Name, sizeCount)
	} else {
		path = fmt.Sprintf("%s/%s (%s)", path, t.Value.Name, sizeCount)
	}

	var v1 float64
	var v2 float64
	if p.ByCount {
		v1, v2 = float64(t.Value.Count), float64(t.Value.Size)
	} else {
		v1, v2 = float64(t.Value.Size), log(t.Value.Count)
	}

	fmt.Fprintf(
		sb,
		"%s,%f,%f\n",
		strings.Replace(path, ",", "-", -1),
		v1,
		v2,
	)

	if p.ByExtension && t.Value.IsDir {
		for _, info := range t.Value.Extensions {
			pth := path + "/" + info.Name
			if p.ByCount {
				v1, v2 = float64(info.Count), float64(info.Size)
			} else {
				v1, v2 = float64(info.Size), log(info.Count)
			}
			fmt.Fprintf(
				sb,
				"%s (%s | %s),%f,%f\n",
				strings.Replace(pth, ",", "-", -1),
				util.FormatUnits(info.Size, "B"),
				util.FormatUnits(int64(info.Count), ""),
				v1,
				v2,
			)
		}
	}
	for _, child := range t.Children {
		if !p.ByExtension || child.Value.IsDir {
			p.print(child, sb, path)
		}
	}
}

func log(n int) float64 {
	return math.Log10(math.Max(float64(n), 1.0))
}
