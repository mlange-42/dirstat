package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mlange-42/dirstat/filesys"
	"github.com/mlange-42/dirstat/tree"
	"github.com/mlange-42/dirstat/util"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dirstat [flags] [command]",
	Short: "Analyze or visualize directory contents.",
	Args:  cobra.NoArgs,
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
	exclude, err := cmd.Flags().GetStringSlice("exclude")
	if err != nil {
		panic(err)
	}
	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		panic(err)
	}
	depth, err := cmd.Flags().GetInt("depth")
	if err != nil {
		panic(err)
	}
	if isJSON && !cmd.Flags().Changed("depth") {
		depth = -1
	}

	var t *tree.FileTree

	if isJSON {
		subtree, serr := cmd.Flags().GetString("subtree")
		if serr != nil {
			panic(serr)
		}
		t, err = treeFromJSON(dir, subtree, exclude, depth)
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
					fmt.Fprintf(os.Stderr, "\rScan: %s, %d files in %s    ", util.FormatUnits(size, "B"), count, time.Since(startTime).Round(time.Millisecond))
				}
			}
		case t = <-done:
			if !quiet {
				fmt.Fprintf(os.Stderr, "\rDone: %s, %d (%s) files in %s    \n", util.FormatUnits(size, "B"), count, util.FormatUnits(int64(count), ""), time.Since(startTime).Round(time.Millisecond))
			}
			break Loop
		case err = <-erro:
			break Loop
		}
	}

	return t, err
}

func treeFromJSON(file string, subtree string, exclude []string, depth int) (*tree.FileTree, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	t, err := tree.Deserialize(bytes)
	if err != nil {
		return nil, err
	}

	if len(subtree) != 0 {
		elems := strings.Split(filepath.ToSlash(subtree), "/")
		t, err = tree.SubTree(t, elems, func(e *tree.FileEntry, path string) bool {
			return strings.ToLower(e.Name) == strings.ToLower(path)
		})
		if err != nil {
			return nil, err
		}
	}

	if depth >= 0 {
		t.Crop(depth, func(parent, child *tree.FileEntry) {
			if child.IsDir {
				parent.AddExtensions(child.Extensions)
			}
		})
	}

	return t, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", ".", "Path to scan or JSON file to load")
	rootCmd.PersistentFlags().StringP("subtree", "s", "", "When reading from JSON, use only this sub-tree")
	rootCmd.PersistentFlags().IntP("depth", "d", 2, "Depth of the generated file tree.\nUse -1 for unlimited depth (use with caution on deeply nested directory trees).\nDefaults to 2 when working on a directory, and to -1 when reading from JSON\n")
	rootCmd.PersistentFlags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns. Ignored when reading from JSON.\nRequires a comma-separated list of patterns, like \"*.exe,.git\"")
	rootCmd.PersistentFlags().Bool("debug", false, "Debug mode with error traces")
	rootCmd.PersistentFlags().Bool("quiet", false, "Don't show progress on stderr")
}
