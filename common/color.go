package common

import (
	"github.com/fatih/color"
	"text/template"
)

var ClorFuncMap = template.FuncMap{
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
}