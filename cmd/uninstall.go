package main

import (
	"fmt"
	"github.com/cnlubo/myssh/myssh"
	"github.com/cnlubo/myssh/utils"
	"github.com/pkg/errors"
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
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
