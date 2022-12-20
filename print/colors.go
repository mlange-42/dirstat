package print

import "github.com/gookit/color"

var (
	directoryColor  = color.C256(39, false).Sprint
	hiddenDirColor  = color.S256(39, 238).Sprint
	fileColor       = color.C256(15, false).Sprint
	hiddenFileColor = color.S256(15, 238).Sprint
	extensionColor  = color.C256(11, false).Sprint
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
