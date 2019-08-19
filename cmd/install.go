package main

import (
	"fmt"
	"github.com/cnlubo/myssh/myssh"
	"github.com/cnlubo/myssh/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"runtime"
)

var installDesc = "Install myssh."

type InstallCommand struct {
	baseCommand
}

// Init initialize command.
func (cc *InstallCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "install",
		Short: "install myssh",
		Long:  installDesc,
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
				return cc.runInstall()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *InstallCommand) addFlags() {
	// TODO // TODO: add flags here
}

func (cc *InstallCommand) runInstall() error {

	var installDir string
	osName := runtime.GOOS
	switch osName {
	case "darwin":
		installDir = "/usr/local/Cellar"
	case "linux":
		installDir = "/opt/modules"
		// default:
		// 	installDir = "/opt/modules"
	}

	return utils.Install(osName, installDir)
}
