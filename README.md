# dirstat

[![Tests](https://github.com/mlange-42/dirstat/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/dirstat/actions/workflows/tests.yml)

A command line tool for analyzing and visualizing disk usage.

## Installation

Using Go:

```shell
go install github.com/mlange-42/dirstat@latest
```

## Usage

Get help:

```shell
dirstat -h
```

### Examples

Run in the current folder, with default settings and JSON output

```shell
dirstat
```

Run in the current folder, with default settings and treemap SVG output piped to a file:

```shell
dirstat treemap > out.svg
```

Statistics over file extensions:

```shell
dirstat treemap -x > out.svg
```

Open the created SVG with the default associated program (ideally a web browser):

```shell
dirstat treemap > out.svg && out.svg
```

Exclude files and directories by glob patterns:

```shell
dirstat -e .git,*.exe
```

Analyze with a different depth than the default of 2:

```shell
dirstat -d 4
```

## References

* Uses [`github.com/nikolaydubina/treemap`](https://github.com/nikolaydubina/treemap) for treemap SVG rendering
