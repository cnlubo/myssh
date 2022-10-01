package myssh

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cnlubo/myssh/version"
	"github.com/fatih/color"
	"time"
)

const (
	MaskPrompt = '*'
)

// error def
var (
	inputEmptyErr   = errors.New("input is empty")
	inputTooLongErr = errors.New("input length must be <= 12")
	notNumberErr    = errors.New("only number support")
)

// GetBakFileName generates a backup dir name by current date and time
func GetBakDirName(name string) string {

	return fmt.Sprintf("%s-%s", name, time.Now().Format("20060102150405"))
}

var logo = `%s

%s V%s
%s

`

func Displaylogo() {

	banner, _ := base64.StdEncoding.DecodeString(version.BannerBase64)
	fgColor := color.New(color.FgHiGreen, color.Bold)
	fgColorGithub := color.New(color.FgMagenta, color.Bold).SprintFunc()
	_, _ = fgColor.Printf(logo, banner, version.Appname, version.Version, fgColorGithub(version.GitHub))
}
