package prompt

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"strconv"
	"strings"
	"text/template"
)

const esc = "\033["
const (
	DefaultErrorMsgPrefix = "✘"
	// DefaultGoodPrompt     = "✔"
	//DefaultPrompt        = ""
	DefaultConfirmPrompt = "?"
)
const DONE = "DONE"

type attribute int

// Foreground weight/decoration attributes.
const (
	reset attribute = iota

	FGBold
	FGFaint
	FGItalic
	FGUnderline
)

// Foreground color attributes
const (
	FGBlack attribute = iota + 30
	FGRed
	FGGreen
	FGYellow
	FGBlue
	FGMagenta
	FGCyan
	FGWhite
)

// Background color attributes
const (
	BGBlack attribute = iota + 40
	BGRed
	BGGreen
	BGYellow
	BGBlue
	BGMagenta
	BGCyan
	BGWhite
)

// ResetCode is the character code used to reset the terminal formatting
var ResetCode = fmt.Sprintf("%s%dm", esc, reset)

const (
	hideCursor       = esc + "?25l"
	showCursor       = esc + "?25h"
	clearLine        = esc + "2K"
	clearDown        = esc + "J"
	clearStartOfLine = esc + "1K"
	clearScreen      = esc + "2J"
	moveUp           = esc + "1F"
	move2Up          = esc + "2A"
	moveDown         = esc + "1B"
	clearTerminal    = "\033c"
)

// FuncMap defines template helpers for the output. It can be extended as a
// regular map.
var FuncMap = template.FuncMap{
	"black":     Styler(FGBlack),
	"red":       Styler(FGRed),
	"green":     Styler(FGGreen),
	"yellow":    Styler(FGYellow),
	"blue":      Styler(FGBlue),
	"magenta":   Styler(FGMagenta),
	"cyan":      Styler(FGCyan),
	"white":     Styler(FGWhite),
	"bgBlack":   Styler(BGBlack),
	"bgRed":     Styler(BGRed),
	"bgGreen":   Styler(BGGreen),
	"bgYellow":  Styler(BGYellow),
	"bgBlue":    Styler(BGBlue),
	"bgMagenta": Styler(BGMagenta),
	"bgCyan":    Styler(BGCyan),
	"bgWhite":   Styler(BGWhite),
	"bold":      Styler(FGBold),
	"faint":     Styler(FGFaint),
	"italic":    Styler(FGItalic),
	"underline": Styler(FGUnderline),
}

func upLine(n uint) string {
	return movementCode(n, 'A')
}

func movementCode(n uint, code rune) string {
	return esc + strconv.FormatUint(uint64(n), 10) + string(code)
}

// Styler returns a func that applies the attributes given in the Styler call
// to the provided string.
func Styler(attrs ...attribute) func(interface{}) string {
	attrstrs := make([]string, len(attrs))
	for i, v := range attrs {
		attrstrs[i] = strconv.Itoa(int(v))
	}

	seq := strings.Join(attrstrs, ";")

	return func(v interface{}) string {
		end := ""
		s, ok := v.(string)
		if !ok || !strings.HasSuffix(s, ResetCode) {
			end = ResetCode
		}
		return fmt.Sprintf("%s%sm%v%s", esc, seq, v, end)
	}
}

var term = termenv.ColorProfile()

// FontColor sets the color of the given string and bolds the font
func FontColor(str, color string) string {
	return termenv.String(str).Foreground(term.Color(color)).Bold().String()
}

// GenMask generate a mask string of the specified length
func GenMask(l int) string {
	return GenStr(l, "*")
}

// GenStr generate a string of the specified length, the string is composed of the given characters
func GenStr(l int, s string) string {
	var ss string
	for i := 0; i < l; i++ {
		ss += s
	}
	return ss
}

func Done() tea.Msg {
	return DONE
}