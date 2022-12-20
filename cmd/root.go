package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/mlange-42/dirstat/filesys"
	"github.com/mlange-42/dirstat/print"
	"github.com/mlange-42/dirstat/tree"
	"github.com/mlange-42/dirstat/util"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dirstat [flags] [command]",
	Short: "Analyze and visualize directory content and disk usage.",
	Long: `Analyze and visualize directory content and disk usage.

When used without a subcommand, prints the result of the analysis as plain-text directory tree.
Analyzes the current directory by default. Use flag --path to analyze a different location.

For graphical visualization, see subcommand 'treemap'.
  $ dirstat treemap -h

To store the result of the analysis for later re-use, see subcommand 'json'.
  $ dirstat json -h
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		byExt, err := cmd.Flags().GetBool("extensions")
		if err != nil {
			panic(err)
		}
		sort, err := cmd.Flags().GetString("sort")
		if err != nil {
			panic(err)
		}
		dirs, err := cmd.Flags().GetBool("dirs")
		if err != nil {
			panic(err)
		}
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			panic(err)
		}
		depth, err := cmd.Flags().GetInt("depth")
		if err != nil {
			panic(err)
		}
		hasDepth := cmd.Flags().Changed("depth")
		if !hasDepth && byExt {
			depth = 0
		}
		colorExp, err := cmd.Flags().GetFloat64("exp")
		if err != nil {
			panic(err)
		}
		if colorExp <= 0 {
			fmt.Fprint(os.Stderr, "ERROR: Color exponent --exp must be greater than 0.0\n")
			os.Exit(1)
		}

		noColors, err := cmd.Flags().GetBool("no-colors")
		if err != nil {
			panic(err)
		}

		if noColors || !color.Support256Color() || !isTerminal() {
			color.Disable()
		}

		if sort != print.ByName && sort != print.BySize && sort != print.ByCount && sort != print.ByAge {
			if debug {
				panic(err)
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: Unknown sort field '%s'. Must be one of [name, size, count, age].\n", sort)
				os.Exit(1)
			}
		}
		t, err := runRootCommand(cmd, args, depth, true)

		if err != nil {
			if debug {
				panic(err)
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				os.Exit(1)
			}
		}

		printer := print.NewFileTreePrinter(byExt, 2, true, dirs, colorExp)
		printer.SortBy = sort
		fmt.Print(printer.Print(t))
	},
}

func runRootCommand(cmd *cobra.Command, args []string, depth int, hasDepth bool) (*tree.FileTree, error) {
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

	isJSON := !info.IsDir() && strings.ToLower(path.Ext(info.Name())) == ".json"
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
	doProfiling, err := cmd.Flags().GetBool("profile")
	if err != nil {
		panic(err)
	}
	if isJSON && !hasDepth {
		depth = -1
	}

	var t *tree.FileTree

	if doProfiling {
		defer profile.Start().Stop()
	}

	if isJSON {
		subtree, serr := cmd.Flags().GetString("select")
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
	progress := make(chan int64, 32)
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
					fmt.Fprintf(os.Stderr, "\rScan: %6s, %d files in %s    ", util.FormatUnits(size, "B"), count, time.Since(startTime).Round(time.Millisecond))
				}
			}
		case t = <-done:
			if !quiet {
				fmt.Fprintf(os.Stderr, "\rDone: %6s, %d (%s) files in %s    \n", util.FormatUnits(size, "B"), count, util.FormatUnitsSimple(int64(count), ""), time.Since(startTime).Round(time.Millisecond))
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

func isTerminal() bool {
	o, _ := os.Stdout.Stat()
	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
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
	rootCmd.PersistentFlags().String("select", "", "When reading from JSON, use only this sub-tree")
	rootCmd.PersistentFlags().StringSliceP("exclude", "e", []string{}, "Exclusion glob patterns. Ignored when reading from JSON.\nRequires a comma-separated list of patterns, like \"*.exe,.git\"")
	rootCmd.PersistentFlags().Bool("debug", false, "Debug mode with error traces")
	rootCmd.PersistentFlags().Bool("quiet", false, "Don't show progress on stderr")
	rootCmd.PersistentFlags().Bool("profile", false, "Do CPU profiling of the analysis part")

	rootCmd.Flags().IntP("depth", "d", 1, "Depth of the generated file tree.\nDeeper files are included, but not individually listed.\nUse -1 for unlimited depth (use with caution on deeply nested directory trees).\nDefaults to -1 when reading from JSON\n")
	rootCmd.Flags().BoolP("extensions", "x", false, "Show directory content by file extension instead of individual files")
	rootCmd.Flags().StringP("sort", "s", "name", "Sort by one of [name, size, count, age]")
	rootCmd.Flags().Bool("dirs", false, "List only directories, no individual files")
	rootCmd.Flags().Float64("exp", 5.0, "Color scale exponent.\n1.0 is linear. Higher values look more log-like.")
	rootCmd.Flags().BoolP("no-colors", "C", false, "Print without colors")
}
