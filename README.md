# dirstat

[![Tests](https://github.com/mlange-42/dirstat/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/dirstat/actions/workflows/tests.yml)

A command line tool for analyzing and visualizing disk usage.

![Screenshot](https://user-images.githubusercontent.com/44003176/208201884-13a4675c-10fa-439f-8b28-21f297a08887.svg)  
*Example visualizing the Go repository*

## Installation

**Using Go:**

```shell
go install github.com/mlange-42/dirstat@latest
```

**Without Go:**

Download binaries for your OS from the [Releases](https://github.com/mlange-42/dirstat/releases/).

## Features

* Visualize disk usage by as text-based tree or as graphical treemap (SVG)
* Optional visualization of directory content by file extension
* Exclusion of files and directories by glob patterns
* Adjustable depth for individual display vs. aggregation
* Write analysis to JSON and re-read for visualization, for handling large directories
* Determines the size of large directories 4x faster than Windows Explorer, and 3x faster than PowerShell

## Usage

Get help:

```shell
dirstat -h
dirstat <command> -h
```

### Basic usage

When used without a subcommand, the result of the analysis is written to STDOUT in JSON format.
If piped to a file, it can be re-used for visualization by using it in the `--path` flag.

Analyze the current directory and write JSON to `out.json`:

```shell
dirstat > out.json
```

Run for a different directory (paths can be absolute or relative):

```shell
dirstat --path ../.. > out.json
```

Read the JSON instead of running an analysis, and print the directory tree in plain text format:

```shell
dirstat --path out.json plain
```

#### Options

Exclude files and directories by glob patterns:

```shell
dirstat --exclude .git,*.exe
```

Analyze with a different depth than the default of 2:

```shell
dirstat --depth 4
```

### Text output

To view a text-based directory tree, use the `plain` command:

```shell
dirstat plain
```

Produces output like this:

```text
dirstat/ ............ 48 kB  (24)    9 minutes
├─.github/ .......... 3.3 kB (2)     4 days
│ └─workflows/ ...... 3.3 kB (2)     4 days
├─cmd/ .............. 12 kB  (3)     9 minutes
│ ├─plain.go ........ 1.4 kB         9 minutes
│ ├─root.go ......... 5.5 kB         15 minutes
│ └─treemap.go ...... 5.3 kB         9 minutes
...
```

#### Options

Aggregate by file extensions and sort by size

```shell
dirstat plain -x --sort size
```

For more options, see the CLI help `dirstat plain -h`.

### Treemap

To generate graphical treemaps, use the `treemap` command.

Generate the treemap and write it to `out.svg` (can be viewed with any web browser):

```shell
dirstat treemap > out.svg
```

Immediately open the created SVG with the default associated program (ideally a web browser):

```shell
dirstat treemap > out.svg && out.svg
```

#### Options

Statistics over file extensions:

```shell
dirstat treemap -x > out.svg
```

Size boxes by file count instead of size:

```shell
dirstat treemap --count > out.svg
```

Produce CSV output for use with [`github.com/nikolaydubina/treemap`](https://github.com/nikolaydubina/treemap):

```shell
dirstat treemap --csv
```

For more options to customize the treemap, see the CLI help `dirstat treemap -h`.

## References

* Uses [`github.com/nikolaydubina/treemap`](https://github.com/nikolaydubina/treemap) for treemap SVG rendering
