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
	Args: cobra.ExactArgs(1),
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

		dir := path.Clean(args[0])
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("%s does not exist", dir))
		}
		isJSON := strings.ToLower(path.Ext(info.Name())) == ".json"
		if !info.IsDir() && !isJSON {
			panic(fmt.Errorf("%s is neither a directory nor a JSON file", dir))
		}

		var t *tree.FileTree

		if isJSON {
			t, err = treeFromJSON(dir, exclude, depth)
		} else {
			t, err = treeFromDir(dir, exclude, depth)
		}
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
	rootCmd.Flags().IntP("depth", "d", 2, "Depth of the file tree.")
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns.")
	rootCmd.Flags().StringP("format", "f", "json", "Output format. One of [plain json].")

	rootCmd.Flags().Bool("plain", false, "Output as plain directory tree.")
	rootCmd.Flags().Bool("treemap", false, "Output as treemap csv.")
	rootCmd.MarkFlagsMutuallyExclusive("plain", "treemap")
}
