package main

import (
	"encoding/base64"
	"github.com/cnlubo/myssh/utils"
	"github.com/cnlubo/myssh/version"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionTpl = `%s

%s V%s
%s

OsArch:   %s_%s
GoVersion:%s
BuildTime:%s
GitCommit:%s
`

var versionDescription = "Display the version information of myssh， " +
	"including GoVersion, KernelVersion, Os, Version, Arch, BuildTime and GitCommit."

// VersionCommand use to implement 'version' command.
type VersionCommand struct {
	baseCommand
}

// Init initialize version command.
func (v *VersionCommand) Init(c *Cli) {
	v.cli = c
	v.cmd = &cobra.Command{
		Use:   "version",
		Short: "Print versions about Myssh",
		Args:  cobra.NoArgs,
		Long:  versionDescription,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return v.runVersion()
		},
	}
	v.addFlags()
}

// addFlags adds flags for specific command.
func (v *VersionCommand) addFlags() {
	// TODO: add flags here
}

// runVersion is the entry of version command.
func (v *VersionCommand) runVersion() error {

	result, err := utils.Version()
	if err != nil {
		utils.ExitN(utils.Err, "failed to get system version:"+err.Error(), 1)
	}
	banner, _ := base64.StdEncoding.DecodeString(version.BannerBase64)
	// fmt.Printf(color.FgLightGreen.Render(versionTpl), banner, version.Appname, version.Version, version.GitHub, result.Os, result.Arch, result.GoVersion, result.BuildTime, result.GitCommit)
	c := color.New(color.FgHiGreen, color.Bold)
	_, _ = c.Printf(versionTpl, banner, version.Appname, version.Version, version.GitHub, result.Os, result.Arch, result.GoVersion, result.BuildTime, result.GitCommit)
	return nil
}