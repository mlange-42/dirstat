package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mlange42/dirstat/crawl"
	"github.com/mlange42/dirstat/tree"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dirstat <path> [flags] command",
	Short: "Analyze or visualize directory contents.",
	Long: `Analyze or visualize directory contents.
Path can be a directory or a JSON file of a previously generated directory tree.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := runRootCommand(cmd, args)

		if err != nil {
			if d, _ := cmd.Flags().GetBool("debug"); d {
				panic(err)
			} else {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		printer := tree.JSONPrinter[*tree.FileEntry]{}
		fmt.Println(printer.Print(t))
	},
}

func runRootCommand(cmd *cobra.Command, args []string) (*tree.FileTree, error) {
	dir, err := cmd.Flags().GetString("path")
	if err != nil {
		panic(err)
	}
	exclude, err := cmd.Flags().GetStringSlice("exclude")
	if err != nil {
		panic(err)
	}
	depth, err := cmd.Flags().GetInt("depth")
	if err != nil {
		panic(err)
	}

	dir = path.Clean(dir)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s does not exist", dir)
		}
		return nil, err
	}
	isJSON := strings.ToLower(path.Ext(info.Name())) == ".json"
	if !info.IsDir() && !isJSON {
		return nil, fmt.Errorf("%s is neither a directory nor a JSON file", dir)
	}

	var t *tree.FileTree

	if isJSON {
		t, err = treeFromJSON(dir, exclude, depth)
	} else {
		t, err = treeFromDir(dir, exclude, depth)
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func treeFromDir(dir string, exclude []string, depth int) (*tree.FileTree, error) {
	t, err := crawl.Walk(dir, exclude, depth)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func treeFromJSON(file string, exclude []string, depth int) (*tree.FileTree, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	t, err := tree.Deserialize(bytes)
	if err != nil {
		return nil, err
	}
	return t, nil
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
	rootCmd.PersistentFlags().StringP("path", "p", ".", "Path to scan or JSON file to load.")
	rootCmd.PersistentFlags().IntP("depth", "d", 2, "Depth of the file tree.")
	rootCmd.PersistentFlags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns.")
	rootCmd.PersistentFlags().Bool("debug", false, "Debug mode with error traces.")
}
