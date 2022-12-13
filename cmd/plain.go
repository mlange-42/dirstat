package cmd

import (
	"fmt"
	"os"

	"github.com/mlange42/dirstat/tree"
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

		printer := tree.PlainPrinter[*tree.FileEntry]{}
		fmt.Print(printer.Print(t))
	},
}

func init() {
	rootCmd.AddCommand(plainCmd)
}
