package print

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mlange-42/dirstat/tree"
	"github.com/mlange-42/dirstat/util"
)

// TreemapPrinter prints a tree in treemap CSV format
type TreemapPrinter struct {
	ByExtension bool
	ByCount     bool
	HeatAge     bool
	OnlyDirs    bool
	currTime    time.Time
}

// NewTreemapPrinter creates a new TreemapPrinter
func NewTreemapPrinter(byExtension, byCount, heatAge, onlyDirs bool) TreemapPrinter {
	return TreemapPrinter{
		ByExtension: byExtension,
		ByCount:     byCount,
		HeatAge:     heatAge,
		OnlyDirs:    onlyDirs,
		currTime:    time.Now(),
	}
}

// Print prints a FileTree
func (p TreemapPrinter) Print(t *tree.FileTree) string {
	sb := strings.Builder{}
	p.print(t, &sb, "")
	return sb.String()
}

func (p TreemapPrinter) print(t *tree.FileTree, sb *strings.Builder, path string) {
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
		if child.Value.IsDir || !(p.ByExtension || p.OnlyDirs) {
			p.print(child, sb, path)
		}
	}
}

func log(n int) float64 {
	return math.Log10(math.Max(float64(n), 1.0))
}
