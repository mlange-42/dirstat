package cmd

import (
	"fmt"
	"os"

	"github.com/mlange42/dirstat/crawl"
	"github.com/mlange42/dirstat/tree"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run path [flags]",
	Short: "Analyze directory content.",
	Long:  `Analyze directory content.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exclude, err := cmd.Flags().GetStringSlice("exclude")
		if err != nil {
			panic(err)
		}
		depth, err := cmd.Flags().GetInt("depth")
		if err != nil {
			panic(err)
		}

		format := "json"
		pl, err := cmd.Flags().GetBool("plain")
		if err != nil {
			panic(err)
		}
		if pl {
			format = "plain"
		}

		tm, err := cmd.Flags().GetBool("treemap")
		if err != nil {
			panic(err)
		}
		if tm {
			format = "treemap"
		}

		t, err := crawl.Walk(args[0], exclude, depth)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var printer tree.Printer[*tree.FileEntry]

		switch format {
		case "plain":
			printer = tree.PlainPrinter[*tree.FileEntry]{}
			fmt.Println(printer.Print(t))
		case "json":
			printer = tree.JSONPrinter[*tree.FileEntry]{}
		case "treemap":
			printer = tree.TreemapPrinter{ByExtension: true}
		}
		fmt.Println(printer.Print(t))
	},
}

func init() {
	runCmd.Flags().IntP("depth", "d", 2, "Depth of the file tree.")
	runCmd.Flags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns.")
	runCmd.Flags().StringP("format", "f", "json", "Output format. One of [plain json].")

	runCmd.Flags().Bool("plain", false, "Output as plain directory tree.")
	runCmd.Flags().Bool("treemap", false, "Output as treemap csv.")
	runCmd.MarkFlagsMutuallyExclusive("plain", "treemap")

	rootCmd.AddCommand(runCmd)
}
