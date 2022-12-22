package print

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

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
	SortBy        string
	Cutoff        float64
	ByExtension   bool
	Indent        int
	PrintTime     bool
	OnlyDirs      bool
	ColorExponent float64
	prefixNone    string
	prefixEmpty   string
	prefixNormal  string
	prefixLast    string
	printWidth    int
	currTime      time.Time
	ageRange      minMax
	countRange    minMax
	sizeRange     minMax
}

// NewFileTreePrinter creates a new FileTreePrinter
func NewFileTreePrinter(byExt bool, cutoff float64, indent int, printTime bool, onlyDirs bool, colorExponent float64) FileTreePrinter {
	return FileTreePrinter{
		ByExtension:   byExt,
		Cutoff:        cutoff,
		Indent:        indent,
		PrintTime:     printTime,
		OnlyDirs:      onlyDirs,
		ColorExponent: colorExponent,
		prefixNone:    strings.Repeat(" ", indent),
		prefixEmpty:   "│" + strings.Repeat(" ", indent-1),
		prefixNormal:  "├" + strings.Repeat("─", indent-1),
		prefixLast:    "└" + strings.Repeat("─", indent-1),
		printWidth:    0,
		currTime:      time.Now(),
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
	pad := strings.Repeat(".", int(math.Max(float64(p.printWidth-depth*p.Indent-strLen(t.Value.Name)), 0)))
	fmt.Fprint(sb, pref)
	if t.Value.IsDir {
		sizeStr := fmt.Sprintf(" %6s ", util.FormatUnits(t.Value.Size, "B"))
		countStr := fmt.Sprintf(" %5s ", util.FormatUnits(int64(t.Value.Count), ""))

		sizeStr = p.sizeRange.Interpolate(float64(t.Value.Size), p.ColorExponent)(sizeStr)
		countStr = p.countRange.Interpolate(float64(t.Value.Count), p.ColorExponent)(countStr)

		nameColor := directoryColor
		if depth > 0 && strings.HasPrefix(t.Value.Name, ".") {
			nameColor = hiddenDirColor
		}
		fmt.Fprintf(sb, "%s %s %s %s", nameColor(t.Value.Name+"/"), pad, sizeStr, countStr)
	} else {
		sizeStr := fmt.Sprintf(" %6s ", util.FormatUnits(t.Value.Size, "B"))

		sizeStr = p.sizeRange.Interpolate(float64(t.Value.Size), p.ColorExponent)(sizeStr)

		nameColor := fileColor
		if depth > 0 && strings.HasPrefix(t.Value.Name, ".") {
			nameColor = hiddenFileColor
		}
		fmt.Fprintf(sb, "%s .%s %s        ", nameColor(t.Value.Name), pad, sizeStr)
	}

	if p.PrintTime {
		val := fmt.Sprintf(" %11s ", util.FormatDuration(t.Value.Time, p.currTime))
		fmt.Fprintf(sb, " %s", p.ageRange.Interpolate(float64(p.currTime.Unix()-t.Value.Time.Unix()), p.ColorExponent)(val))
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
		sorter := FileEntrySorter{children, func(t *tree.FileTree) float64 { return float64(t.Value.Size) }}
		children = sorter.Sort(p.Cutoff)
	case ByCount:
		sorter := FileEntrySorter{children, func(t *tree.FileTree) float64 { return float64(t.Value.Count) }}
		children = sorter.Sort(p.Cutoff)
	case ByAge:
		sorter := FileEntrySorter{children, func(t *tree.FileTree) float64 { return -float64(t.Value.Time.Unix()) }}
		children = sorter.Sort(1.0)
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
	values := maps.Values(ext)
	switch p.SortBy {
	case BySize:
		sorter := ExtensionEntrySorter{values, func(e *tree.ExtensionEntry) float64 { return float64(e.Size) }}
		values = sorter.Sort(p.Cutoff)
	case ByCount:
		sorter := ExtensionEntrySorter{values, func(e *tree.ExtensionEntry) float64 { return float64(e.Count) }}
		values = sorter.Sort(p.Cutoff)
	case ByAge:
		sorter := ExtensionEntrySorter{values, func(e *tree.ExtensionEntry) float64 { return -float64(e.Time.Unix()) }}
		values = sorter.Sort(1.0)
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
		pad := strings.Repeat(".", int(math.Max(float64(p.printWidth-(depth)*p.Indent-strLen(info.Name)), 0)))
		if i == len(values)-1 {
			fmt.Fprint(sb, prefLast)
		} else {
			fmt.Fprint(sb, pref)
		}

		sizeStr := fmt.Sprintf(" %6s ", util.FormatUnits(info.Size, "B"))
		countStr := fmt.Sprintf(" %5s ", util.FormatUnits(int64(info.Count), ""))

		sizeStr = p.sizeRange.Interpolate(float64(info.Size), p.ColorExponent)(sizeStr)
		countStr = p.countRange.Interpolate(float64(info.Count), p.ColorExponent)(countStr)
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
			fmt.Fprintf(sb, " %s", p.ageRange.Interpolate(float64(p.currTime.Unix()-info.Time.Unix()), p.ColorExponent)(val))
		}
		fmt.Fprint(sb, "\n")

	}
}

func (p FileTreePrinter) maxWidth(t *tree.FileTree, depth int, extensions bool) int {
	max := strLen(t.Value.Name) + depth*p.Indent
	if extensions && t.Value.IsDir {
		for name := range t.Value.Extensions {
			m := strLen(name) + (depth+1)*p.Indent
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
	p.calcAgeRange(t, p.ByExtension)
	p.calcSizeRange(t, p.ByExtension)
	p.calcCountRange(t, p.ByExtension)
}

func (p *FileTreePrinter) calcAgeRange(t *tree.FileTree, extensions bool) {
	unix := p.currTime.Unix()
	p.ageRange.min, p.ageRange.max, _ = p.calcRange(t, extensions, true,
		func(e *tree.FileEntry) (float64, bool) {
			if e.Time.IsZero() {
				return 1, false
			}
			return math.Max(1, float64(unix-e.Time.Unix())), true
		},
		func(e *tree.ExtensionEntry) (float64, bool) {
			if e.Time.IsZero() {
				return 1, false
			}
			return math.Max(1, float64(unix-e.Time.Unix())), true
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

// FileEntrySorter sorts file entries
type FileEntrySorter struct {
	Slice  []*tree.FileTree
	Getter func(*tree.FileTree) float64
}

func (p FileEntrySorter) Len() int      { return len(p.Slice) }
func (p FileEntrySorter) Swap(i, j int) { p.Slice[i], p.Slice[j] = p.Slice[j], p.Slice[i] }
func (p FileEntrySorter) Less(i, j int) bool {
	return p.Getter(p.Slice[i]) > p.Getter(p.Slice[j])
}

// Sort sorts with a cutoff
func (p FileEntrySorter) Sort(cutoff float64) []*tree.FileTree {
	sort.Sort(p)

	if cutoff >= 1.0 {
		return p.Slice
	}

	total := 0.0
	for _, e := range p.Slice {
		total += p.Getter(e)
	}
	max := total * cutoff
	isCut := false

	total = 0.0
	result := []*tree.FileTree{}
	remainder := tree.NewDir("<cutoff>")
	skipped := 0
	for _, e := range p.Slice {
		value := p.Getter(e)
		if isCut {
			remainder.Value.Add(e.Value.Size, e.Value.Count, e.Value.Time)
			skipped++
		} else {
			result = append(result, e)
		}
		total += value
		if !isCut && total > max {
			isCut = true
		}
	}

	if skipped > 0 {
		remainder.Value.Name = fmt.Sprintf("<skipped %d>", skipped)
		result = append(result, remainder)
	}

	return result
}

// ExtensionEntrySorter sorts file entries
type ExtensionEntrySorter struct {
	Slice  []*tree.ExtensionEntry
	Getter func(*tree.ExtensionEntry) float64
}

func (p ExtensionEntrySorter) Len() int      { return len(p.Slice) }
func (p ExtensionEntrySorter) Swap(i, j int) { p.Slice[i], p.Slice[j] = p.Slice[j], p.Slice[i] }
func (p ExtensionEntrySorter) Less(i, j int) bool {
	return p.Getter(p.Slice[i]) > p.Getter(p.Slice[j])
}

// Sort sorts with a cutoff
func (p ExtensionEntrySorter) Sort(cutoff float64) []*tree.ExtensionEntry {
	sort.Sort(p)

	if cutoff >= 1.0 {
		return p.Slice
	}

	total := 0.0
	for _, e := range p.Slice {
		total += p.Getter(e)
	}
	max := total * cutoff
	isCut := false

	total = 0.0
	result := []*tree.ExtensionEntry{}
	remainder := tree.ExtensionEntry{Name: "<cutoff>"}
	skipped := 0
	for _, e := range p.Slice {
		value := p.Getter(e)
		if isCut {
			remainder.Add(e.Size, e.Count, e.Time)
			skipped++
		} else {
			result = append(result, e)
		}
		total += value
		if !isCut && total > max {
			isCut = true
		}
	}

	if skipped > 0 {
		remainder.Name = fmt.Sprintf("<skipped %d>", skipped)
		result = append(result, &remainder)
	}

	return result
}

type minMax struct {
	min float64
	max float64
}

func (r minMax) Interpolate(value float64, exponent float64) func(a ...interface{}) string {
	if r.min >= r.max {
		return defaultColors[0]
	}

	rel := math.Pow(math.Max(0, value), 1.0/exponent) / math.Pow(math.Max(0, r.max), 1.0/exponent)

	index := int(rel * float64(len(defaultColors)+1))
	if index < 0 {
		index = 0
	}
	if index >= len(defaultColors) {
		index = len(defaultColors) - 1
	}
	return defaultColors[index]
}

func strLen(str string) int {
	return utf8.RuneCountInString(str)
}
