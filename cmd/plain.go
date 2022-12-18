package cmd

import (
	"fmt"
	"os"

	"github.com/mlange-42/dirstat/tree"
	"github.com/spf13/cobra"
)

// treemapCmd represents the treemap command
var plainCmd = &cobra.Command{
	Use:   "plain",
	Short: "Prints a plain text directory tree",
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
		sort, err := cmd.Flags().GetString("sort")
		if err != nil {
			panic(err)
		}
		if sort != "" && sort != tree.BySize && sort != tree.ByCount {
			fmt.Fprintf(os.Stderr, "Unknown sort field '%s'. Must be one of [size, count]", sort)
			os.Exit(1)
		}

		printer := tree.NewFileTreePrinter(byExt, 2)
		printer.SortBy = sort
		fmt.Print(printer.Print(t))
	},
}

func init() {
	plainCmd.Flags().BoolP("extensions", "x", false, "Show directory content by file extension instead of individual files")
	plainCmd.Flags().String("sort", "", "Sort by either size or count. Possible values: [size, count]")

	rootCmd.AddCommand(plainCmd)
}
