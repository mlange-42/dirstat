package cmd

import (
	"fmt"
	"os"

	"github.com/mlange-42/dirstat/print"
	"github.com/spf13/cobra"
)

// treemapCmd represents the treemap command
var plainCmd = &cobra.Command{
	Use:   "plain",
	Short: "Prints a plain text directory tree.",
	Run: func(cmd *cobra.Command, args []string) {
		byExt, err := cmd.Flags().GetBool("extensions")
		if err != nil {
			panic(err)
		}
		sort, err := cmd.Flags().GetString("sort")
		if err != nil {
			panic(err)
		}
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			panic(err)
		}

		if sort != print.ByName && sort != print.BySize && sort != print.ByCount && sort != print.ByAge {
			if debug {
				panic(err)
			} else {
				fmt.Fprintf(os.Stderr, "Unknown sort field '%s'. Must be one of [name, size, count, age].\n", sort)
				os.Exit(1)
			}
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

		printer := print.NewFileTreePrinter(byExt, 2, true)
		printer.SortBy = sort
		fmt.Print(printer.Print(t))
	},
}

func init() {
	plainCmd.Flags().BoolP("extensions", "x", false, "Show directory content by file extension instead of individual files")
	plainCmd.Flags().String("sort", "name", "Sort by one of [name, size, count, age]")

	rootCmd.AddCommand(plainCmd)
}
