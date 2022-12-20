package print

import "github.com/gookit/color"

var (
	red       = color.FgRed.Render
	green     = color.FgGreen.Render
	lightBlue = color.FgLightBlue.Render
	yellow    = color.FgYellow.Render
)

var defaultColors = []func(a ...interface{}) string{
	color.C256(17, true).Sprint,
	color.C256(17, true).Sprint,
	color.C256(18, true).Sprint,
	color.C256(18, true).Sprint,
	color.C256(23, true).Sprint,
	color.C256(22, true).Sprint,
	color.C256(58, true).Sprint,
	color.C256(94, true).Sprint,
	color.C256(88, true).Sprint,
	color.C256(52, true).Sprint,
}

var colorMaps = map[string][]func(a ...interface{}) string{
	"default": defaultColors,
}
