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

		printer := tree.FileTreePrinter{ByExtension: byExt}
		fmt.Print(printer.Print(t))
	},
}

func init() {
	plainCmd.Flags().BoolP("extensions", "x", false, "Show directory content by file extension instead of individual files")

	rootCmd.AddCommand(plainCmd)
}
