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
