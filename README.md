# go-dirstat

A command line tool for analyzing and visualizing disk usage.

## Installation

Using Go:

```shell
go install github.com/mlange-42/go-dirstat@latest
```

## Usage

Get help:

```shell
dirstat -h
```

Run in the current folder, with default settings and JSON output

```shell
dirstat
```

Run in the current folder, with default settings and SVG output piped to a file:

```shell
dirstat treemap --svg > out.svg
```

Statistics over file extensions:

```shell
dirstat treemap --svg -x > out.svg
```

## References

* Uses [`github.com/nikolaydubina/treemap`](https://github.com/nikolaydubina/treemap) for treemap SVG rendering
