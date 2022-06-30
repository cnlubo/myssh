package common

import (
	"github.com/fatih/color"
	"text/template"
)

var ColorFuncMap = template.FuncMap{
	"black":     color.New(color.FgHiBlack).SprintFunc(),
	"red":       color.New(color.FgHiRed).SprintFunc(),
	"green":     color.New(color.FgHiGreen).SprintFunc(),
	"yellow":    color.New(color.FgHiYellow).SprintFunc(),
	"blue":      color.New(color.FgHiBlue).SprintFunc(),
	"magenta":   color.New(color.FgHiMagenta).SprintFunc(),
	"cyan":      color.New(color.FgHiCyan).SprintFunc(),
	"white":     color.New(color.FgHiWhite).SprintFunc(),
	"bgBlack":   color.New(color.BgHiBlack).SprintFunc(),
	"bgRed":     color.New(color.BgHiRed).SprintFunc(),
	"bgGreen":   color.New(color.BgHiGreen).SprintFunc(),
	"bgYellow":  color.New(color.BgHiYellow).SprintFunc(),
	"bgBlue":    color.New(color.BgHiBlue).SprintFunc(),
	"bgMagenta": color.New(color.BgHiMagenta).SprintFunc(),
	"bgCyan":    color.New(color.BgHiCyan).SprintFunc(),
	"bgWhite":   color.New(color.BgHiWhite).SprintFunc(),
	"bold":      color.New(color.Bold).SprintFunc(),
	"faint":     color.New(color.Faint).SprintFunc(),
	"italic":    color.New(color.Italic).SprintFunc(),
	"underline": color.New(color.Underline).SprintFunc(),
	"reset":     color.New(color.Reset).SprintFunc(),
}

func RenderedText(str, strcolor string) string {
	switch strcolor {
	case "blue":
		return color.New(color.Bold, color.FgHiBlue).SprintfFunc()(str)
	case "white":
		return color.New(color.Bold, color.FgHiWhite).SprintFunc()(str)
	case "fwhite":
		return color.New(color.Faint, color.FgHiWhite).SprintFunc()(str)
	case "red":
		return color.New(color.Bold, color.FgHiRed).SprintFunc()(str)
	case "cyan":
		return color.New(color.Bold, color.FgHiCyan).SprintFunc()(str)
	case "green":
		return color.New(color.Bold, color.FgHiGreen).SprintFunc()(str)
	case "yellow":
		return color.New(color.Bold, color.FgHiYellow).SprintFunc()(str)
	}
	return str
}