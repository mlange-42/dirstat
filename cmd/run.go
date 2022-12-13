/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mlange42/dirstat/crawl"
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
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			panic(err)
		}
		if _, ok := map[string]bool{"plain": true, "json": true}[format]; !ok {
			fmt.Printf("Format option '%s' unknown. Must be one of [plain json]\n", format)
			os.Exit(1)
		}

		t, err := crawl.Walk(args[0], exclude, depth)
		if err != nil {
			panic(err)
		}

		switch format {
		case "plain":
			fmt.Println(t)
		case "json":
			tt, err := json.MarshalIndent(t, "", "    ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(tt[:]))
			t := crawl.NewDir("root")
			err = json.Unmarshal(tt, &t)
			if err != nil {
				panic(err)
			}
			fmt.Println(t)
		}
	},
}

func init() {
	runCmd.Flags().IntP("depth", "d", 2, "Depth of the file tree.")
	runCmd.Flags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns.")
	runCmd.Flags().StringP("format", "f", "json", "Output format. One of [plain json].")

	rootCmd.AddCommand(runCmd)
}
