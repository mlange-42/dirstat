package print

import "github.com/gookit/color"

var (
	red       = color.FgRed.Render
	green     = color.FgGreen.Render
	lightBlue = color.FgLightBlue.Render
	yellow    = color.FgYellow.Render
)

var defaultColors = []func(a ...interface{}) string{
	color.FgBlue.Render,
	color.FgLightBlue.Render,
	color.FgCyan.Render,
	color.FgGreen.Render,
	color.FgYellow.Render,
	color.FgRed.Render,
}

var colorMaps = map[string][]func(a ...interface{}) string{
	"default": defaultColors,
}
