package cmd

import (
	"fmt"
	"os"

	"github.com/mlange42/dirstat/crawl"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dirstat path [flags]",
	Short: "Shows statistics about directory contents.",
	Long:  `Shows statistics about directory contents.`,
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
		tree, err := crawl.Walk(args[0], exclude, depth)
		if err != nil {
			panic(err)
		}
		fmt.Println(tree)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("depth", "d", 2, "Depth of the file tree.")
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns.")
}
