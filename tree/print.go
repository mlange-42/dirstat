package tree

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/mlange-42/dirstat/util"
	"golang.org/x/exp/maps"
)

const (
	// BySize is for sorting by size
	BySize string = "size"
	// ByCount is for sorting by count
	ByCount string = "count"
	// ByAge is for sorting by age
	ByAge string = "age"
	// ByName is for sorting by name
	ByName string = "name"
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
	SortBy       string
	ByExtension  bool
	Indent       int
	PrintTime    bool
	prefixNone   string
	prefixEmpty  string
	prefixNormal string
	prefixLast   string
	printWidth   int
	currTime     time.Time
}

// NewFileTreePrinter creates a new FileTreePrinter
func NewFileTreePrinter(byExt bool, indent int, printTime bool) FileTreePrinter {
	return FileTreePrinter{
		ByExtension:  byExt,
		Indent:       indent,
		PrintTime:    printTime,
		prefixNone:   strings.Repeat(" ", indent),
		prefixEmpty:  "│" + strings.Repeat(" ", indent-1),
		prefixNormal: "├" + strings.Repeat("─", indent-1),
		prefixLast:   "└" + strings.Repeat("─", indent-1),
		printWidth:   0,
		currTime:     time.Now(),
	}
}

// Print prints a FileTree
func (p FileTreePrinter) Print(t *FileTree) string {
	p.printWidth = p.maxWidth(t, 0, p.ByExtension) + 1
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
		sizeCount = fmt.Sprintf("%-6s (%s)",
			util.FormatUnits(t.Value.Size, "B"), util.FormatUnits(int64(t.Value.Count), ""),
		)
	} else {
		sizeCount = fmt.Sprintf("%s", util.FormatUnits(t.Value.Size, "B"))
	}

	pref := prefix

	if depth > 0 {
		pref = prefix + p.createPrefix(last)
	}
	pad := strings.Repeat(".", int(math.Max(float64(p.printWidth-depth*p.Indent-len([]rune(t.Value.Name))), 0)))
	fmt.Fprint(sb, pref)
	if t.Value.IsDir {
		fmt.Fprintf(sb, "%s/ %s %-15s", t.Value.Name, pad, sizeCount)
	} else {
		fmt.Fprintf(sb, "%s .%s %-15s", t.Value.Name, pad, sizeCount)
	}

	if p.PrintTime {
		util.FPrintDuration(sb, t.Value.Time, p.currTime)
	}
	fmt.Fprint(sb, "\n")

	if depth > 0 {
		pref = prefix + p.createPrefixEmpty(last)
	}

	children := t.Children
	switch p.SortBy {
	case BySize:
		sorter := SortDesc[FileTree]{children, func(t *FileTree) float64 { return float64(t.Value.Size) }}
		sort.Sort(sorter)
	case ByCount:
		sorter := SortDesc[FileTree]{children, func(t *FileTree) float64 { return float64(t.Value.Count) }}
		sort.Sort(sorter)
	case ByAge:
		sorter := SortDesc[FileTree]{children, func(t *FileTree) float64 { return -float64(t.Value.Time.Unix()) }}
		sort.Sort(sorter)
	case ByName:
	default:
		panic(fmt.Errorf("Unknown sort field '%s'", p.SortBy))
	}

	for i, child := range children {
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

	values := maps.Values(ext)
	switch p.SortBy {
	case BySize:
		sorter := SortDesc[ExtensionEntry]{values, func(e *ExtensionEntry) float64 { return float64(e.Size) }}
		sort.Sort(sorter)
	case ByCount:
		sorter := SortDesc[ExtensionEntry]{values, func(e *ExtensionEntry) float64 { return float64(e.Count) }}
		sort.Sort(sorter)
	case ByAge:
		sorter := SortDesc[ExtensionEntry]{values, func(e *ExtensionEntry) float64 { return -float64(e.Time.Unix()) }}
		sort.Sort(sorter)
	case ByName:
		sort.Slice(values, func(i, j int) bool {
			return values[i].Name < values[j].Name
		})
	default:
		panic(fmt.Errorf("Unknown sort field '%s'", p.SortBy))
	}

	pref := prefix + p.createPrefix(false)
	prefLast := prefix + p.createPrefix(true)

	for i, info := range values {
		pad := strings.Repeat(".", int(math.Max(float64(p.printWidth-(depth)*p.Indent-len([]rune(info.Name))), 0)))
		if i == len(keys)-1 {
			fmt.Fprint(sb, prefLast)
		} else {
			fmt.Fprint(sb, pref)
		}
		sizeCount := fmt.Sprintf("%-6s (%s)", util.FormatUnits(info.Size, "B"), util.FormatUnits(int64(info.Count), ""))
		fmt.Fprintf(
			sb,
			"%s .%s %-15s",
			info.Name,
			pad,
			sizeCount,
		)

		if p.PrintTime {
			util.FPrintDuration(sb, info.Time, p.currTime)
		}
		fmt.Fprint(sb, "\n")

	}
}

func (p FileTreePrinter) maxWidth(t *FileTree, depth int, extensions bool) int {
	max := len([]rune(t.Value.Name)) + depth*p.Indent
	if t.Value.IsDir {
		for name := range t.Value.Extensions {
			m := len([]rune(name)) + (depth+1)*p.Indent
			if m > max {
				max = m
			}
		}
	}
	for _, c := range t.Children {
		v := p.maxWidth(c, depth+1, extensions)
		if v > max {
			max = v
		}
	}
	return max
}

func (p FileTreePrinter) createPrefix(last bool) string {
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
	HeatAge     bool
	currTime    time.Time
}

// NewTreemapPrinter creates a new TreemapPrinter
func NewTreemapPrinter(byExtension bool, byCount bool, heatAge bool) TreemapPrinter {
	return TreemapPrinter{
		ByExtension: byExtension,
		ByCount:     byCount,
		HeatAge:     heatAge,
		currTime:    time.Now(),
	}
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
	if p.HeatAge {
		v2 = p.currTime.Sub(t.Value.Time).Hours() / 24
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
			if p.HeatAge {
				v2 = p.currTime.Sub(info.Time).Hours() / 24
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

// SortDesc sorts by size
type SortDesc[T any] struct {
	Slice  []*T
	Getter func(*T) float64
}

func (p SortDesc[T]) Len() int      { return len(p.Slice) }
func (p SortDesc[T]) Swap(i, j int) { p.Slice[i], p.Slice[j] = p.Slice[j], p.Slice[i] }
func (p SortDesc[T]) Less(i, j int) bool {
	return p.Getter(p.Slice[i]) > p.Getter(p.Slice[j])
}
