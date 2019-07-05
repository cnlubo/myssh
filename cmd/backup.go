package main

import (
	"github.com/cnlubo/myssh/myssh"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"path/filepath"
)

var BackupDesc = "Backup all config、SSHKeys."

type BackupCommand struct {
	baseCommand
	backupPath string
}

func (cc *BackupCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup all config,SSHKeys...",
		Long:  BackupDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runBackup()

		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *BackupCommand) addFlags() {
	flagSet := cc.cmd.Flags()
	home, _ := homedir.Dir()
	backPath := filepath.Join(home, ".mysshbackup")
	flagSet.StringVarP(&cc.backupPath, "path", "p", backPath, "backup path")

}

func (cc *BackupCommand) runBackup() error {

	return myssh.BackupAll(cc.backupPath, &cc.cli.Env)

}
