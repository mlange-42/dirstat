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
		exclude := map[string]bool{
			".git": true,
		}
		tree, err := crawl.Walk(args[0], exclude)
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
}
