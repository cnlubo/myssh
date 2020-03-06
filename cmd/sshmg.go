package main

import (
	"fmt"
	"github.com/cnlubo/myssh/myssh"
	"github.com/cnlubo/myssh/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
	"sort"
	"strings"
)

var hostAliasDesc = "command line tool for managing your ssh alias config."

type HostAliasCommand struct {
	baseCommand
}

// Init initialize command.
func (cc *HostAliasCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "alias",
		Aliases: []string{"malias"},
		Short:   "managing your ssh alias config (alias: malias)",
		Long:    hostAliasDesc,
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
				myssh.Displaylogo()
				return cc.cmd.Help()
			}

		},
	}

	c.AddCommand(cc, &aliasListCmd{})
	c.AddCommand(cc, &aliasDeleteCmd{})
	c.AddCommand(cc, &aliasAddCmd{})
	c.AddCommand(cc, &aliasUpdateCmd{})
	c.AddCommand(cc, &aliasBatchCmd{})
	c.AddCommand(cc, &aliasGoCmd{})
	c.AddCommand(cc, &AliasKeyCopyCmd{})
}

var aliasListDesc = "List all ssh alias."

var aliasListHeader = []string{
	" AliasName ",
	" ConnetString ",
	" Config ",
}

type aliasListCmd struct {
	baseCommand
	ignoreCase bool
	showPath   bool
}

func (cc *aliasListCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "list Keywords...",
		Aliases: []string{"ls"},
		Short:   "List ssh alias (alias:ls)",
		Long:    aliasListDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runAliasList(args)
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *aliasListCmd) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.SetInterspersed(false)
	flagSet.BoolVarP(&cc.ignoreCase, "ignore", "i", true, "ignore case while searching")
	flagSet.BoolVarP(&cc.showPath, "path", "p", false, "display the file path of the alias")

}

func (cc *aliasListCmd) runAliasList(args []string) error {

	sshConfigPath := filepath.Join(cc.cli.Env.SSHPath, "config")
	hosts, err := myssh.ListAlias(sshConfigPath, myssh.ListOption{
		Keywords:   args,
		IgnoreCase: cc.ignoreCase,
	})
	if err != nil {
		return errors.Wrapf(err, "get Alias List failed")
	}
	if len(hosts) > 0 {
		utils.PrintN(utils.Info, fmt.Sprintf("ssh alias total records: %d\n", len(hosts)))
		fmt.Println()
		cc.cli.PrintTable(displayAliases(cc.showPath, hosts))
	} else {
		utils.PrintN(utils.Info, "not found ssh alias\n")
	}
	return nil
}

var deleteAliasDesc = " Delete one or more ssh aliases."

type aliasDeleteCmd struct {
	baseCommand
}

func (cc *aliasDeleteCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "delete alias...",
		Aliases: []string{"del"},
		Short:   "Delete ssh alias (alias:del)",
		Long:    deleteAliasDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, -1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runAliasDelete(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *aliasDeleteCmd) addFlags() {
	// TODO add Flags here

}
func (cc *aliasDeleteCmd) runAliasDelete(args []string) error {
	sshConfigPath := filepath.Join(cc.cli.Env.SSHPath, "config")
	_, err := myssh.DeleteAliases(sshConfigPath, args...)
	if err != nil {
		return errors.Wrap(err, "deleted failed")
	} else {
		utils.PrintN(utils.Info, "deleted successfully!!!\n\n")
		// cc.cli.PrintTable(displayAliases(false, hosts))
	}
	return nil
}

var addAliasDesc = "Add a new SSH alias record."

type aliasAddCmd struct {
	baseCommand
	identityfile string
}

func (cc *aliasAddCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new SSH alias record",
		Long:  addAliasDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runAliasAdd()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *aliasAddCmd) addFlags() {

	// TODO add Flags here
}
func (cc *aliasAddCmd) runAliasAdd() error {

	sshConfigPath := filepath.Join(cc.cli.Env.SSHPath, "config")
	hosts, err := myssh.AddHostAlias(sshConfigPath)
	if err != nil {
		return err
	} else {
		cc.cli.PrintTable(displayAliases(false, hosts))
	}
	return nil
}

var updateAliasDesc = "Update the specified ssh alias."

type aliasUpdateCmd struct {
	baseCommand
	identityfile string
}

func (cc *aliasUpdateCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "update alias",
		Short: "Update ssh alias",
		Long:  updateAliasDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runAliasUpdate(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *aliasUpdateCmd) addFlags() {

	// TODO add Flags here
}
func (cc *aliasUpdateCmd) runAliasUpdate(args []string) error {

	sshConfigPath := filepath.Join(cc.cli.Env.SSHPath, "config")
	hosts, err := myssh.UpdateHostCfg(sshConfigPath, args[0])
	if err != nil {
		return err
	} else {
		cc.cli.PrintTable(displayAliases(false, hosts))
	}
	return nil
}

var batchAliasDesc = "Batch exec command for alias."

type aliasBatchCmd struct {
	baseCommand
	Prompt     bool
	ignoreCase bool
}

// Init initialize command.
func (cc *aliasBatchCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "batch alias command ... [flags]",
		Aliases: []string{"bt"},
		Short:   "batch exec command (alias: bt)",
		Long:    batchAliasDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 2, -1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runExec(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *aliasBatchCmd) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.BoolVarP(&cc.Prompt, "prompt", "P", false, "Prompt for password")
}

func (cc *aliasBatchCmd) runExec(args []string) error {

	sshConfigPath := filepath.Join(cc.cli.Env.SSHPath, "config")
	return myssh.BatchExecAlias(sshConfigPath, cc.Prompt, args[0], args[1:]...)

}

var AliasKeyCopyDesc = "Copy SSH public key to alias Host"

type AliasKeyCopyCmd struct {
	baseCommand
	identityfile string
}

func (cc *AliasKeyCopyCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "keycopy alias...",
		Aliases: []string{"kcp"},
		Short:   "Copy SSH public key to a alias Host (alias:kcp)",
		Long:    AliasKeyCopyDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, -1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runAliasKeyCopy(args)
			}
		},
	}
	cc.addFlags()
}

func (cc *AliasKeyCopyCmd) addFlags() {

	// TODO: add flags here
}
func (cc *AliasKeyCopyCmd) runAliasKeyCopy(args []string) error {

	sshConfigPath := filepath.Join(cc.cli.Env.SSHPath, "config")
	return myssh.AliasKeyCopy(sshConfigPath, args)
}

var AliasGoDesc = "ssh login Server."

type aliasGoCmd struct {
	baseCommand
	promptPass bool
}

func (cc *aliasGoCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "go alias",
		Short: "ssh login Server",
		Long:  AliasGoDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runAliasGo(args)
			}

		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *aliasGoCmd) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.BoolVarP(&cc.promptPass, "prompt", "P", false, "Prompt for password")
}

func (cc *aliasGoCmd) runAliasGo(args []string) error {

	sshConfigFile := filepath.Join(cc.cli.Env.SSHPath, "config")
	fmt.Printf(sshConfigFile)
	if len(args) == 0 {
		return myssh.AliasInteractiveLogin(sshConfigFile, cc.promptPass)
	} else {
		return myssh.AliasLogin(args[0], cc.promptPass, sshConfigFile)
	}

}

func displayAliases(showPath bool, hosts []*myssh.HostConfig) (tableHeader []string, hColors []tablewriter.Colors, colColors []tablewriter.Colors, data [][]string) {

	var aliases []string
	var noConnectAliases []string

	hostMap := map[string]*myssh.HostConfig{}

	for _, host := range hosts {

		hostMap[host.Alias] = host

		if host.Display() {
			aliases = append(aliases, host.Alias)
		} else {
			noConnectAliases = append(noConnectAliases, host.Alias)
		}

	}

	var rowData [][]string
	sort.Strings(aliases)
	for _, alias := range aliases {
		aliasRow := addRow(showPath, hostMap[alias])
		rowData = append(rowData, aliasRow)
	}

	sort.Strings(noConnectAliases)
	for _, noConnectAlias := range noConnectAliases {
		noConnectAliasRow := addRow(showPath, hostMap[noConnectAlias])
		rowData = append(rowData, noConnectAliasRow)
	}

	var headerColors, columnColors []tablewriter.Colors
	for i := 1; i <= len(aliasListHeader); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold})
		columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgGreenColor})
	}

	return aliasListHeader, headerColors, columnColors, rowData

}

func addRow(showPath bool, host *myssh.HostConfig) (row []string) {
	var hostConfig []string
	aliasName := host.Alias
	if showPath && len(host.PathMap) > 0 {
		var paths []string
		for path := range host.PathMap {
			homeDir, _ := homedir.Dir()
			if strings.HasPrefix(path, homeDir) {
				path = strings.Replace(path, homeDir, "~", 1)
			}
			paths = append(paths, path)
		}
		sort.Strings(paths)
		aliasName = aliasName + "(" + strings.Join(paths, " ") + ")"
	}
	connectstr := host.ConnectionStr()
	if connectstr == "" {
		connectstr = "-"
	}
	var config string
	var configs []string
	for _, key := range utils.SortKeys(host.OwnConfig) {
		value := host.OwnConfig[key]
		if value == "" {
			continue
		}
		if key == "identityfile" {
			homeDir, _ := homedir.Dir()
			if strings.HasPrefix(value, homeDir) {
				value = strings.Replace(value, homeDir, "~", 1)
			}
		}
		configs = append(configs, key+":"+value)
	}

	for _, key := range utils.SortKeys(host.ImplicitConfig) {
		value := host.ImplicitConfig[key]
		if value == "" {
			continue
		}
		configs = append(configs, key+":"+value)
	}
	config = strings.Join(configs, "\n")
	if len(config) == 0 {
		config = "-"
	}
	hostConfig = append(hostConfig, aliasName, connectstr, config)
	return hostConfig
}
