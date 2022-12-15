package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mlange-42/dirstat/filesys"
	"github.com/mlange-42/dirstat/tree"
	"github.com/mlange-42/dirstat/util"
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
	quiet, err := cmd.Flags().GetBool("quiet")
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
		t, err = treeFromDir(dir, exclude, depth, quiet)
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func treeFromDir(dir string, exclude []string, depth int, quiet bool) (*tree.FileTree, error) {
	progress := make(chan int64)
	done := make(chan *tree.Tree[*tree.FileEntry])
	erro := make(chan error)

	var t *tree.FileTree = nil
	var err error = nil
	var size int64 = 0
	var count int = 0
	minElapsed := 250 * time.Millisecond

	go filesys.Walk(dir, exclude, depth, progress, done, erro)

	startTime := time.Now()
	prevTime := startTime

Loop:
	for {
		select {
		case p := <-progress:
			size += p
			count++
			if !quiet {
				if count%10 == 0 && time.Since(prevTime) >= minElapsed {
					prevTime = time.Now()
					fmt.Fprintf(os.Stderr, "\033[2K\rScanning: %s, %d files in %s", util.FormatUnits(size, "B"), count, time.Since(startTime).Round(time.Millisecond))
				}
			}
		case t = <-done:
			if !quiet {
				fmt.Fprintf(os.Stderr, "\033[2K\rDone: %s, %s files in %s", util.FormatUnits(size, "B"), util.FormatUnits(int64(count), ""), time.Since(startTime).Round(time.Millisecond))
			}
			break Loop
		case err = <-erro:
			break Loop
		}
	}

	return t, err
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
	t.Crop(depth)
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
	rootCmd.PersistentFlags().Bool("quiet", false, "Don't show progress on stderr.")
}
