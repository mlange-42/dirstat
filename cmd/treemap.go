package cmd

import (
	"fmt"
	"image/color"
	"os"

	"github.com/mlange-42/dirstat/tree"
	"github.com/nikolaydubina/treemap"
	"github.com/nikolaydubina/treemap/parser"
	"github.com/nikolaydubina/treemap/render"
	"github.com/spf13/cobra"
)

// treemapCmd represents the treemap command
var treemapCmd = &cobra.Command{
	Use:     "treemap",
	Aliases: []string{"tm"},
	Short:   "Visualize output as SVG treemap.",
	Run: func(cmd *cobra.Command, args []string) {
		byExt, err := cmd.Flags().GetBool("extensions")
		if err != nil {
			panic(err)
		}

		byCount, err := cmd.Flags().GetBool("count")
		if err != nil {
			panic(err)
		}

		csv, err := cmd.Flags().GetBool("csv")
		if err != nil {
			panic(err)
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			panic(err)
		}

		t, err := runRootCommand(cmd, args)
		if err != nil {
			if debug {
				panic(err)
			} else {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		printer := tree.TreemapPrinter{
			ByExtension: byExt,
			ByCount:     byCount,
		}
		str := printer.Print(t)
		if csv {
			fmt.Print(str)
			return
		}

		svgFlags := parseSvgFlags(cmd)

		svgBytes, err := toSvg(str, &svgFlags)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(svgBytes)
	},
}

var grey = color.RGBA{128, 128, 128, 255}

func toSvg(s string, flags *svgFlags) ([]byte, error) {
	parser := parser.CSVTreeParser{}
	tree, err := parser.ParseString(s)
	if err != nil || tree == nil {
		return []byte{}, err
	}

	treemap.SetNamesFromPaths(tree)
	if !flags.KeepLongPaths {
		treemap.CollapseLongPaths(tree)
	}

	sizeImputer := treemap.SumSizeImputer{EmptyLeafSize: 1}
	sizeImputer.ImputeSize(*tree)

	if flags.ImputeHeat {
		heatImputer := treemap.WeightedHeatImputer{EmptyLeafHeat: 0.5}
		heatImputer.ImputeHeat(*tree)
	}

	tree.NormalizeHeat()

	var colorer render.Colorer

	palette, hasPalette := render.GetPalette(flags.ColorScheme)
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
	case flags.ColorScheme == "none":
		colorer = render.NoneColorer{}
		borderColor = grey
	case flags.ColorScheme == "balanced":
		colorer = treeHueColorer
		borderColor = color.White
	case hasPalette && tree.HasHeat():
		colorer = render.HeatColorer{Palette: palette}
	case tree.HasHeat():
		palette, _ := render.GetPalette("RdBu")
		colorer = render.HeatColorer{Palette: palette}
	default:
		colorer = treeHueColorer
	}

	switch {
	case flags.ColorBorder == "light":
		borderColor = color.White
	case flags.ColorBorder == "dark":
		borderColor = grey
	}

	uiBuilder := render.UITreeMapBuilder{
		Colorer:     colorer,
		BorderColor: borderColor,
	}
	spec := uiBuilder.NewUITreeMap(*tree, flags.W, flags.H, flags.MarginBox, flags.PaddingBox, flags.Padding)
	renderer := render.SVGRenderer{}

	return renderer.Render(spec, flags.W, flags.H), nil
}

func init() {
	treemapCmd.Flags().Bool("csv", false, "Generate raw CSV output for github.com/nikolaydubina/treemap")
	treemapCmd.Flags().BoolP("extensions", "x", false, "Show directory content by file extension instead of individual files")
	treemapCmd.Flags().BoolP("count", "c", false, "Size boxes by file count instead of disk memory")

	treemapCmd.Flags().Float64("w", 1028, "width of output")
	treemapCmd.Flags().Float64("h", 640, "height of output")
	treemapCmd.Flags().Float64("margin-box", 4, "margin between boxes")
	treemapCmd.Flags().Float64("padding-box", 4, "padding between box border and content")
	treemapCmd.Flags().Float64("padding", 32, "padding around root content")
	treemapCmd.Flags().String("color", "balance", "color scheme (RdBu, balance, RdYlGn, none)")
	treemapCmd.Flags().String("color-border", "auto", "color of borders (light, dark, auto)")
	treemapCmd.Flags().Bool("impute-heat", false, "impute heat for parents(weighted sum) and leafs(0.5)")
	treemapCmd.Flags().Bool("long-paths", false, "keep long paths when paren has single child")

	rootCmd.AddCommand(treemapCmd)
}

type svgFlags struct {
	W             float64
	H             float64
	MarginBox     float64
	PaddingBox    float64
	Padding       float64
	ColorScheme   string
	ColorBorder   string
	ImputeHeat    bool
	KeepLongPaths bool
}

func parseSvgFlags(cmd *cobra.Command) svgFlags {
	var flags svgFlags

	flags.W, _ = cmd.Flags().GetFloat64("w")
	flags.H, _ = cmd.Flags().GetFloat64("h")
	flags.MarginBox, _ = cmd.Flags().GetFloat64("margin-box")
	flags.PaddingBox, _ = cmd.Flags().GetFloat64("padding-box")
	flags.Padding, _ = cmd.Flags().GetFloat64("padding")

	flags.ColorScheme, _ = cmd.Flags().GetString("color")
	flags.ColorBorder, _ = cmd.Flags().GetString("color-border")

	flags.ImputeHeat, _ = cmd.Flags().GetBool("impute-heat")
	flags.KeepLongPaths, _ = cmd.Flags().GetBool("long-paths")

	return flags
}
