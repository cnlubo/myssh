package main

import (
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"github.com/spf13/cobra"
	"path/filepath"
	"runtime"
)

var uninstallDescription = "uninstall myssh."

type UninstallCommand struct {
	baseCommand
}

// Init initialize command.
func (cc *UninstallCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall myssh",
		Long:  uninstallDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("command 'myssh uninstall %s' does not exist.\nPlease execute `myssh uninstall --help` for more help", args[0])
			} else {
				return cc.runUninstall()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *UninstallCommand) addFlags() {

	// TODO // TODO: add flags here
}

func (cc *UninstallCommand) runUninstall() error {

	var uninstallDir string

	osName := runtime.GOOS
	switch osName {
	case "darwin":
		uninstallDir = "/usr/local/Cellar"
	case "linux":
		uninstallDir = "/opt/modules"
		// default:
		// 	uninstallDir = "/usr/local/Cellar"
	}

	return utils.Uninstall(osName, filepath.Join(uninstallDir, "myssh"))
}
