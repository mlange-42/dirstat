/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mlange42/dirstat/crawl"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [flags]",
	Short: "Visualize directory content from analysis.",
	Long:  `Visualize directory content from analysis.`,
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		t := crawl.NewDir("root")
		err = json.Unmarshal(bytes, &t)
		if err != nil {
			panic(err)
		}
		fmt.Println(t)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
