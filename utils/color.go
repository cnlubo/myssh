package utils

import (
	"github.com/gookit/color"
	"sort"
	"sync"
	"text/template"
)

const (
	ColorRed     = "red"
	ColorGreen   = "green"
	ColorYellow  = "yellow"
	ColorBlue    = "blue"
	ColorMagenta = "magenta"
	ColorCyan    = "cyan"
	ColorWhite   = "white"
	ColorGray    = "gray"
)

type colorCount struct {
	name  string
	color color.Color
	count int
}

type colorCounts []colorCount

func (cs colorCounts) Len() int {
	return len(cs)
}
func (cs colorCounts) Less(i, j int) bool {
	return cs[i].count < cs[j].count
}
func (cs colorCounts) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

var colorMux sync.Mutex

func GetColor() (string, color.Color) {
	colorMux.Lock()
	defer colorMux.Unlock()
	sort.Sort(cs)
	cs[0].count++
	return cs[0].name, cs[0].color
}

var cs = colorCounts{

	colorCount{ColorRed, color.FgRed, 0},
	colorCount{ColorGreen, color.FgGreen, 0},
	colorCount{ColorYellow, color.FgYellow, 0},
	colorCount{ColorMagenta, color.FgMagenta, 0},
	colorCount{ColorCyan, color.FgCyan, 0},
	colorCount{ColorGray, color.FgGray, 0},
	colorCount{ColorBlue, color.FgBlue, 0},
	// colorCount{ColorWhite, color.FgWhite, 0},
}

var ColorFuncMap = template.FuncMap{

	ColorRed:     color.Style{color.FgRed, color.OpBold}.Render,
	ColorGreen:   color.Style{color.FgGreen, color.OpBold}.Render,
	ColorYellow:  color.Style{color.FgYellow, color.OpBold}.Render,
	ColorMagenta: color.Style{color.FgMagenta, color.OpBold}.Render,
	ColorCyan:    color.Style{color.FgCyan, color.OpBold}.Render,
	ColorGray:    color.Style{color.FgGray, color.OpBold}.Render,
	ColorBlue:    color.Style{color.FgBlue, color.OpBold}.Render,
	// ColorWhite:   color.Style{color.FgCyan, color.OpBold}.Render,
}
