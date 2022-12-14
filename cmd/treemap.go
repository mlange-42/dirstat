package cmd

import (
	"fmt"
	"image/color"
	"os"

	"github.com/mlange42/dirstat/tree"
	"github.com/nikolaydubina/treemap"
	"github.com/nikolaydubina/treemap/parser"
	"github.com/nikolaydubina/treemap/render"
	"github.com/spf13/cobra"
)

// treemapCmd represents the treemap command
var treemapCmd = &cobra.Command{
	Use:   "treemap",
	Short: "Create output in treemap CSV format",
	Run: func(cmd *cobra.Command, args []string) {
		t, err := runRootCommand(cmd, args)
		if err != nil {
			if d, _ := cmd.Flags().GetBool("debug"); d {
				panic(err)
			} else {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		byExt, err := cmd.Flags().GetBool("extensions")
		if err != nil {
			panic(err)
		}

		byCount, err := cmd.Flags().GetBool("count")
		if err != nil {
			panic(err)
		}

		svg, err := cmd.Flags().GetBool("svg")
		if err != nil {
			panic(err)
		}

		printer := tree.TreemapPrinter{
			ByExtension: byExt,
			ByCount:     byCount,
		}
		str := printer.Print(t)
		if !svg {
			fmt.Print(str)
			return
		}

		svgBytes, err := toSvg(str)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(svgBytes)
	},
}

func toSvg(s string) ([]byte, error) {
	parser := parser.CSVTreeParser{}
	tree, err := parser.ParseString(s)
	if err != nil || tree == nil {
		return []byte{}, err
	}

	treemap.SetNamesFromPaths(tree)
	//treemap.CollapseLongPaths(tree)

	sizeImputer := treemap.SumSizeImputer{EmptyLeafSize: 1}
	sizeImputer.ImputeSize(*tree)

	heatImputer := treemap.WeightedHeatImputer{EmptyLeafHeat: 0.5}
	heatImputer.ImputeHeat(*tree)

	tree.NormalizeHeat()

	var colorer render.Colorer

	palette, hasPalette := render.GetPalette("balanced")
	treeHueColorer := render.TreeHueColorer{
		Offset: 0,
		Hues:   map[string]float64{},
		C:      0.5,
		L:      0.5,
		DeltaH: 10,
		DeltaC: 0.3,
		DeltaL: 0.1,
	}

	var borderColor color.Color
	borderColor = color.White
	colorer = treeHueColorer
	borderColor = color.White

	colorer = treeHueColorer
	borderColor = color.White

	switch {
	case hasPalette && tree.HasHeat():
		colorer = render.HeatColorer{Palette: palette}
	case tree.HasHeat():
		palette, _ := render.GetPalette("RdBu")
		colorer = render.HeatColorer{Palette: palette}
	}

	uiBuilder := render.UITreeMapBuilder{
		Colorer:     colorer,
		BorderColor: borderColor,
	}
	spec := uiBuilder.NewUITreeMap(*tree, 1028, 640, 4, 4, 32)
	renderer := render.SVGRenderer{}

	return renderer.Render(spec, 1028, 640), nil
}

func init() {
	treemapCmd.Flags().BoolP("extensions", "x", false, "Group deepest directory level by file extensions.")
	treemapCmd.Flags().BoolP("count", "c", false, "Size boxes by file count instead of disk memory.")
	treemapCmd.Flags().Bool("svg", false, "Directly greates SVG output with default treemap settings.")

	rootCmd.AddCommand(treemapCmd)
}
