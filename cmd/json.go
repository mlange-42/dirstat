package cmd

import (
	"fmt"
	"os"

	"github.com/mlange-42/dirstat/print"
	"github.com/mlange-42/dirstat/tree"
	"github.com/spf13/cobra"
)

// treemapCmd represents the treemap command
var jsonCmd = &cobra.Command{
	Use:     "json",
	Aliases: []string{"js"},
	Short:   "Prints the tree as JSON for later re-use.",
	Long: `Prints the tree as JSON for later re-use.

Writes the result of the analysis to STDOUT in JSON format.
When piped to a file, it can be used for visualization later by passing as the '--path' argument.

  $ dirstat json > out.json
    (analyzes the current directory and writes JSON to out.json)

  $ dirstat --path out.json
    (reads the JSON instead of running an analysis, and prints the directory tree in text format)
`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			panic(err)
		}

		depth, err := cmd.Flags().GetInt("depth")
		if err != nil {
			panic(err)
		}
		hasDepth := cmd.Flags().Changed("depth")

		t, err := runRootCommand(cmd, args, depth, hasDepth)
		if err != nil {
			if debug {
				panic(err)
			} else {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		printer := print.JSONPrinter[*tree.FileEntry]{}
		fmt.Print(printer.Print(t))
	},
}

func init() {
	jsonCmd.Flags().IntP("depth", "d", 2, "Depth of the generated file tree.\nDeeper files are included, but not individually listed.\nUse -1 for unlimited depth (use with caution on deeply nested directory trees).\nDefaults to -1 when reading from JSON\n")

	rootCmd.AddCommand(jsonCmd)
}
