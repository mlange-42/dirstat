package cmd

import (
	"fmt"
	"os"

	"github.com/mlange42/dirstat/tree"
	"github.com/spf13/cobra"
)

// treemapCmd represents the treemap command
var treemapCmd = &cobra.Command{
	Use:   "treemap",
	Short: "Create output in treemap CSV format",
	Run: func(cmd *cobra.Command, args []string) {
		t, err := runRootCommand(cmd, args)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		byExt, err := cmd.Flags().GetBool("extensions")
		if err != nil {
			panic(err)
		}
		printer := tree.TreemapPrinter{ByExtension: byExt}
		fmt.Print(printer.Print(t))
	},
}

func init() {
	treemapCmd.Flags().BoolP("extensions", "x", false, "Group deepest directory level by file extensions.")

	rootCmd.AddCommand(treemapCmd)
}
