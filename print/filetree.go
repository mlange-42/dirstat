package print

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/mlange-42/dirstat/tree"
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

// FileTreePrinter prints a file tree in plain text format
type FileTreePrinter struct {
	SortBy       string
	ByExtension  bool
	Indent       int
	PrintTime    bool
	OnlyDirs     bool
	prefixNone   string
	prefixEmpty  string
	prefixNormal string
	prefixLast   string
	printWidth   int
	currTime     time.Time
	timeRange    minMax
	countRange   minMax
	sizeRange    minMax
}

// NewFileTreePrinter creates a new FileTreePrinter
func NewFileTreePrinter(byExt bool, indent int, printTime bool, onlyDirs bool) FileTreePrinter {
	return FileTreePrinter{
		ByExtension:  byExt,
		Indent:       indent,
		PrintTime:    printTime,
		OnlyDirs:     onlyDirs,
		prefixNone:   strings.Repeat(" ", indent),
		prefixEmpty:  "│" + strings.Repeat(" ", indent-1),
		prefixNormal: "├" + strings.Repeat("─", indent-1),
		prefixLast:   "└" + strings.Repeat("─", indent-1),
		printWidth:   0,
		currTime:     time.Now(),
	}
}

// Print prints a FileTree
func (p FileTreePrinter) Print(t *tree.FileTree) string {
	p.calcRanges(t)

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

func (p FileTreePrinter) print(t *tree.FileTree, sb *strings.Builder, depth int, last bool, prefix string) {
	pref := prefix

	if depth > 0 {
		pref = prefix + p.createPrefix(last)
	}
	pad := strings.Repeat(".", int(math.Max(float64(p.printWidth-depth*p.Indent-len([]rune(t.Value.Name))), 0)))
	fmt.Fprint(sb, pref)
	if t.Value.IsDir {
		sizeStr := fmt.Sprintf(" %6s ", util.FormatUnits(t.Value.Size, "B"))
		countStr := fmt.Sprintf(" %5s ", util.FormatUnits(int64(t.Value.Count), ""))

		sizeStr = p.sizeRange.Interpolate(float64(t.Value.Size), false)(sizeStr)
		countStr = p.countRange.Interpolate(float64(t.Value.Count), false)(countStr)

		nameColor := directoryColor
		if depth > 0 && strings.HasPrefix(t.Value.Name, ".") {
			nameColor = hiddenDirColor
		}
		fmt.Fprintf(sb, "%s %s %s %s", nameColor(t.Value.Name+"/"), pad, sizeStr, countStr)
	} else {
		sizeStr := fmt.Sprintf(" %6s ", util.FormatUnits(t.Value.Size, "B"))

		sizeStr = p.sizeRange.Interpolate(float64(t.Value.Size), false)(sizeStr)

		nameColor := fileColor
		if depth > 0 && strings.HasPrefix(t.Value.Name, ".") {
			nameColor = hiddenFileColor
		}
		fmt.Fprintf(sb, "%s .%s %s        ", nameColor(t.Value.Name), pad, sizeStr)
	}

	if p.PrintTime {
		val := fmt.Sprintf(" %11s ", util.FormatDuration(t.Value.Time, p.currTime))
		fmt.Fprintf(sb, " %s", p.timeRange.Interpolate(float64(t.Value.Time.Unix()), true)(val))
	}
	fmt.Fprint(sb, "\n")

	if depth > 0 {
		pref = prefix + p.createPrefixEmpty(last)
	}

	var children []*tree.FileTree
	for _, child := range t.Children {
		if child.Value.IsDir || !(p.OnlyDirs || p.ByExtension) {
			children = append(children, child)
		}
	}
	switch p.SortBy {
	case BySize:
		sorter := SortDesc[tree.FileTree]{children, func(t *tree.FileTree) float64 { return float64(t.Value.Size) }}
		sort.Sort(sorter)
	case ByCount:
		sorter := SortDesc[tree.FileTree]{children, func(t *tree.FileTree) float64 { return float64(t.Value.Count) }}
		sort.Sort(sorter)
	case ByAge:
		sorter := SortDesc[tree.FileTree]{children, func(t *tree.FileTree) float64 { return -float64(t.Value.Time.Unix()) }}
		sort.Sort(sorter)
	case ByName:
	default:
		panic(fmt.Errorf("Unknown sort field '%s'", p.SortBy))
	}

	for i, child := range children {
		last := i == len(children)-1 && (!p.ByExtension || len(t.Value.Extensions) == 0)
		p.print(child, sb, depth+1, last, pref)
	}

	if p.ByExtension && t.Value.IsDir {
		p.printExtensions(t.Value.Extensions, sb, depth+1, pref)
	}
}

func (p FileTreePrinter) printExtensions(ext map[string]*tree.ExtensionEntry, sb *strings.Builder, depth int, prefix string) {
	keys := make([]string, 0, len(ext))
	for k := range ext {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := maps.Values(ext)
	switch p.SortBy {
	case BySize:
		sorter := SortDesc[tree.ExtensionEntry]{values, func(e *tree.ExtensionEntry) float64 { return float64(e.Size) }}
		sort.Sort(sorter)
	case ByCount:
		sorter := SortDesc[tree.ExtensionEntry]{values, func(e *tree.ExtensionEntry) float64 { return float64(e.Count) }}
		sort.Sort(sorter)
	case ByAge:
		sorter := SortDesc[tree.ExtensionEntry]{values, func(e *tree.ExtensionEntry) float64 { return -float64(e.Time.Unix()) }}
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

		sizeStr := fmt.Sprintf(" %6s ", util.FormatUnits(info.Size, "B"))
		countStr := fmt.Sprintf(" %5s ", util.FormatUnits(int64(info.Count), ""))

		sizeStr = p.sizeRange.Interpolate(float64(info.Size), false)(sizeStr)
		countStr = p.countRange.Interpolate(float64(info.Count), false)(countStr)
		fmt.Fprintf(
			sb,
			"%s .%s %s %s",
			extensionColor(info.Name),
			pad,
			sizeStr,
			countStr,
		)

		if p.PrintTime {
			val := fmt.Sprintf(" %11s ", util.FormatDuration(info.Time, p.currTime))
			fmt.Fprintf(sb, " %s", p.timeRange.Interpolate(float64(info.Time.Unix()), true)(val))
		}
		fmt.Fprint(sb, "\n")

	}
}

func (p FileTreePrinter) maxWidth(t *tree.FileTree, depth int, extensions bool) int {
	max := len([]rune(t.Value.Name)) + depth*p.Indent
	if extensions && t.Value.IsDir {
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

func (p *FileTreePrinter) calcRanges(t *tree.FileTree) {
	p.calcTimeRange(t, p.ByExtension)
	p.calcSizeRange(t, p.ByExtension)
	p.calcCountRange(t, p.ByExtension)
}

func (p *FileTreePrinter) calcTimeRange(t *tree.FileTree, extensions bool) {
	p.timeRange.min, p.timeRange.max, _ = p.calcRange(t, extensions, true,
		func(e *tree.FileEntry) (float64, bool) {
			if e.Time.IsZero() {
				return 0, false
			}
			return float64(e.Time.Unix()), true
		},
		func(e *tree.ExtensionEntry) (float64, bool) {
			if e.Time.IsZero() {
				return 0, false
			}
			return float64(e.Time.Unix()), true
		})
}

func (p *FileTreePrinter) calcSizeRange(t *tree.FileTree, extensions bool) {
	p.sizeRange.min, p.sizeRange.max, _ = p.calcRange(t, extensions, true,
		func(e *tree.FileEntry) (float64, bool) {
			return float64(e.Size), true
		},
		func(e *tree.ExtensionEntry) (float64, bool) {
			return float64(e.Size), true
		})
}

func (p *FileTreePrinter) calcCountRange(t *tree.FileTree, extensions bool) {
	p.countRange.min, p.countRange.max, _ = p.calcRange(t, extensions, true,
		func(e *tree.FileEntry) (float64, bool) {
			return float64(e.Count), true
		},
		func(e *tree.ExtensionEntry) (float64, bool) {
			return float64(e.Count), true
		})
}

func (p FileTreePrinter) calcRange(t *tree.FileTree, extensions bool, isRoot bool,
	fileFn func(*tree.FileEntry) (value float64, on bool),
	extFn func(*tree.ExtensionEntry) (value float64, on bool)) (min float64, max float64, isOk bool) {

	isOk = false

	if isRoot {
		min = math.MaxFloat64
		max = -math.MaxFloat64
		isOk = true
	} else {
		if v, ok := fileFn(t.Value); ok {
			min = v
			max = v
			isOk = true
		} else {
			min = math.MaxFloat64
			max = -math.MaxFloat64
		}
	}

	if extensions && t.Value.IsDir {
		for _, ext := range t.Value.Extensions {
			if v, ok := extFn(ext); ok {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				isOk = true
			}
		}
	}

	for _, c := range t.Children {
		if mn, mx, ok := p.calcRange(c, extensions, false, fileFn, extFn); ok {
			if mn < min {
				min = mn
			}
			if mx > max {
				max = mx
			}
			isOk = true
		}
	}

	return
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

type minMax struct {
	min float64
	max float64
}

func (r minMax) Interpolate(value float64, inverse bool) func(a ...interface{}) string {
	if r.min >= r.max {
		return defaultColors[0]
	}
	rel := (value - r.min) / (r.max - r.min)
	if inverse {
		rel = 1.0 - rel
	}
	index := int(rel * float64(len(defaultColors)+1))
	if index < 0 {
		index = 0
	}
	if index >= len(defaultColors) {
		index = len(defaultColors) - 1
	}
	return defaultColors[index]
}
