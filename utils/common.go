package utils

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const (

	// CheckSymbol is the code for check symbol
	CheckSymbol = "\u2714 "
	// CrossSymbol is the code for check symbol
	CrossSymbol  = "\u2716 "
	ExclamSymbol = "\u0021 "
	ArrowSymbol  = "\u279c"
	DeleteSymbol = "\u2620"
	NoneSymbol   = "\u2605"
	// Message type
	Info     = "info"
	Err      = "error"
	Warn     = "warn"
	Inst     = "install"
	Uinst    = "uninstall"
	None     = "none"
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed, color.Bold).SprintFunc()
	green   = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	magenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
	blue    = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	white   = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	gray    = color.New(color.FgCyan, color.Bold).SprintFunc()
)

func CheckAndExit(err error) {
	if err != nil {
		ExitN(Err, err.Error(), 1)
	}
}
func CheckErr(err error) bool {
	PrintErr(err)
	return err == nil
}
func PrintErr(err error) {
	a := yellow
	if err != nil {
		fmt.Printf("%s%s\n", a(CrossSymbol), red(err.Error()))

	}
}
func PrintErrWithPrefix(prefix string, err error) {
	if err != nil {
		fmt.Println(prefix, err.Error())
	}
}

func PrintN(messageType string, message string) {

	if strings.TrimSpace(messageType) == "" {
		messageType = Info
	}
	if strings.TrimSpace(message) == "" {
		message = "No message"
	}
	switch messageType {
	case Info:
		//fmt.Printf("%s%s\n", color.FgYellow.Render(CheckSymbol), color.FgGreen.Render(message))
		fmt.Printf("%s%s\n", yellow(CheckSymbol), green(message))
	case Err:
		fmt.Printf("%s%s\n", yellow(CrossSymbol), red(message))
	case Warn:
		fmt.Printf("%s%s\n", yellow(ExclamSymbol), magenta(message))
	case Inst:
		fmt.Printf("%s%s\n", yellow(ArrowSymbol), blue(message))
	case Uinst:
		fmt.Printf("%s%s\n", yellow(DeleteSymbol), white(message))
	case None:
		fmt.Printf("%s%s\n", yellow(NoneSymbol), gray(message))

	}
}

func ExitN(messageType string, message string, code int) {

	if strings.TrimSpace(messageType) == "" {
		messageType = Info
	}
	if strings.TrimSpace(message) == "" {
		message = "No message"
	}

	switch messageType {
	case Info:
		fmt.Printf("%s%s\n", yellow(CheckSymbol), green(message))
	case Err:
		fmt.Printf("%s%s\n", red(CrossSymbol), red(message))
	case Warn:
		fmt.Printf("%s%s\n", yellow(ExclamSymbol), magenta(message))

	}
	os.Exit(code)
}

// Execute executes shell commands with arguments
func Execute(workDir, script string, args ...string) bool {

	cmd := exec.Command(script, args...)

	if workDir != "" {
		cmd.Dir = workDir
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		PrintErr(err)
		return false
	}

	return true
}

// Query values contains keys
func Query(values, keys []string, ignoreCase bool) bool {
	contains := func(key string) bool {
		if ignoreCase {
			key = strings.ToLower(key)
		}
		for _, value := range values {
			if ignoreCase {
				value = strings.ToLower(value)
			}
			if strings.Contains(value, key) {
				return true
			}
		}
		return false
	}
	for _, key := range keys {
		if contains(key) {
			return true
		}
	}
	return false
}

// SortKeys sort map keys
func SortKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ArgumentsCheck check arguments count correctness
func ArgumentsCheck(argCount, min, max int) error {
	var err error
	if min > 0 && argCount < min {
		err = errors.New("too few arguments")
	}
	if (max > 0 && argCount > max) || (max == 0 && argCount > 0) {
		err = errors.New("too many arguments")
	}
	return err
}

// check connect string, format is [user@]host[:port]
func CheckConnect(connect string) error {
	var u, h, p string
	hs := strings.SplitN(connect, "@", 2)
	if len(hs) == 2 {
		u = hs[0]
		h = hs[1]
	} else {
		return errors.New(fmt.Sprintf("%s invalid", connect))
	}
	if u == "" {
		return errors.New(fmt.Sprintf("[%s] ==> User is empty", connect))
	}
	hss := strings.SplitN(h, ":", 2)
	if len(hss) == 2 {
		p = hss[1]
		h = hss[0]
		if _, err := strconv.Atoi(p); err != nil {
			return errors.New(fmt.Sprintf("[%s] ==> Port invalid", connect))
		}
	} else {
		return errors.New(fmt.Sprintf("[%s] ==> Port is empty", connect))
	}
	return nil
}

// compress string by regexp trim space or tab
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	reg := regexp.MustCompile("^\\s+|\\s+$")
	return reg.ReplaceAllString(str, "")
}

// IsNumeric is_numeric()
// Numeric strings consist of optional sign, any number of digits, optional decimal part and optional exponential part.
// Thus +0123.45e6 is a valid numeric value.
// In PHP hexadecimal (e.g. 0xf4c3b00c) is not supported, but IsNumeric is supported.
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.TrimSpace(str)
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

func DeleteExtraSpace(s string) string {
	// 删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)   // 替换tab为空格
	regstr := "\\s{2,}"                         // 两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)            // 编译正则表达式
	s2 := make([]byte, len(s1))                 // 定义字符数组切片
	copy(s2, s1)                                // 将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(s2)) // 在字符串中搜索
	for len(spcIndex) > 0 {                     // 找到适配项
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...) // 删除多余空格
		spcIndex = reg.FindStringIndex(string(s2))           // 继续在字符串中搜索
	}
	return string(s2)
}

// 错误断言
func ErrorAssert(err error, assert string) bool {
	return strings.Contains(err.Error(), assert)
}

// clear screen
func Clear() {
	var cmd exec.Cmd
	if "windows" == runtime.GOOS {
		cmd = *exec.Command("cmd", "/c", "cls")
	} else {
		cmd = *exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// ParseConnect parse connect string, format is [user@]host[:port]
func ParseConnect(connectStr string) (string, string, string) {
	var u, hostname, port string
	hs := strings.SplitN(connectStr, "@", 2)
	hostname = hs[0]
	if len(hs) == 2 {
		u = hs[0]
		hostname = hs[1]
	}
	hss := strings.SplitN(hostname, ":", 2)
	hostname = hss[0]
	if len(hss) == 2 {
		if _, err := strconv.Atoi(hss[1]); err == nil {
			port = hss[1]
		}
	}
	return u, hostname, port
}

func Render(tpl *template.Template, data interface{}) []byte {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		return []byte(fmt.Sprintf("%v", data))
	}
	return buf.Bytes()
}

// ErrEOF is the error returned from prompts when EOF is encountered.
var ErrEOF = errors.New("^D")

// ErrInterrupt is the error returned from prompts when an interrupt (ctrl-c) is
// encountered.
var ErrInterrupt = errors.New("^C")

// ErrAbort is the error returned when confirm prompts are supplied "n"
var ErrAbort = errors.New("")