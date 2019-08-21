package main

import (
	"fmt"
	"github.com/cnlubo/myssh/myssh"
	"github.com/cnlubo/myssh/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path"
	"sort"
)

var configDesc = "manage myssh ConfigFile."

type CfgCommand struct {
	baseCommand
}

// Init initialize command.
func (cc *CfgCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "cfg",
		Aliases: []string{"mcfg"},
		Short:   "manage myssh configfile (alias: mcfg)",
		Long:    configDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				// return myssh.InteractiveSetContext(&cc.cli.Env)
				return cc.cmd.Help()
			}
		},
	}
	c.AddCommand(cc, &cfgListCommand{})
	c.AddCommand(cc, &cfgAddCommand{})
	c.AddCommand(cc, &cfgDelCommand{})
	c.AddCommand(cc, &cfgSetCommand{})
}

var listConfigDesc = "List All config."

type cfgListCommand struct {
	baseCommand
	showPath bool
}

// Init initializes command.
func (cc *cfgListCommand) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List config (alias:ls)",
		Long:    listConfigDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runListConfig()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *cfgListCommand) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.SetInterspersed(false)
	flagSet.BoolVarP(&cc.showPath, "path", "p", false, "display the file path")
}

var tableHead = []string{
	"  Name  ",
	"ClusterConfig",
	" SSHConfig ",
}

func (cc *cfgListCommand) runListConfig() error {

	configs := myssh.Main.Contexts
	if configs.Len() == 0 {
		return errors.New("No context found")
	}
	cc.cli.PrintTable(displayConfig(configs, cc.showPath))
	return nil
}

func displayConfig(cfgs myssh.Contexts, showPath bool) (tableHeader []string, hColors []tablewriter.Colors, colColors []tablewriter.Colors, data [][]string) {

	var cfgList []struct {
		myssh.Context
		IsCurrent bool
	}

	sort.Sort(cfgs)
	for _, c := range cfgs {
		cfgList = append(cfgList, struct {
			myssh.Context
			IsCurrent bool
		}{
			Context:   c,
			IsCurrent: c.Name == myssh.Main.Current})
	}

	fmt.Println()

	var rowData [][]string
	var clusterCfg string
	var sshCfg string
	for _, cf := range cfgList {
		var row []string
		if showPath {
			clusterCfg = cf.ClusterConfig
			sshCfg = cf.SSHConfig
		} else {
			clusterCfg = path.Base(cf.ClusterConfig)
			sshCfg = path.Base(cf.SSHConfig)
		}
		if cf.IsCurrent {
			row = append(row, utils.CheckSymbol+cf.Name, clusterCfg, sshCfg)
		} else {
			row = append(row, cf.Name, clusterCfg, sshCfg)
		}
		rowData = append(rowData, row)
	}

	var headerColors, columnColors []tablewriter.Colors
	for i := 1; i <= len(tableHead); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold})
		columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgGreenColor})

	}
	return tableHead, headerColors, columnColors, rowData

}

var delConfigDesc = "delete config."

type cfgDelCommand struct {
	baseCommand
}

// Init initializes command.
func (cc *cfgDelCommand) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "delete cfg...",
		Aliases: []string{"del"},
		Short:   "delete config",
		Long:    delConfigDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, -1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runDelConfig(args)
			}
		},
	}

	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *cfgDelCommand) addFlags() {
	// TODO: add flags here
}

func (cc *cfgDelCommand) runDelConfig(args []string) error {

	if len(myssh.Main.Contexts) == 1 {
		return errors.New("Only one context not allow delete !!!")
	}

	_, err := myssh.DeleteConfig(args, &cc.cli.Env)
	if err != nil {
		return err
	}
	return nil
}

var addConfigDesc = "Add context config."

type cfgAddCommand struct {
	baseCommand
}

// Init initializes command.
func (cc *cfgAddCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "add",
		Short: "Add context config",
		Long:  addConfigDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runAddConfig()
			}
		},
	}

	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *cfgAddCommand) addFlags() {
	// TODO: add flags here
}

func (cc *cfgAddCommand) runAddConfig() error {

	cfgs, err := myssh.AddConfig(&cc.cli.Env)
	if err != nil {
		return err
	}
	cc.cli.PrintTable(displayConfig(cfgs, true))
	return nil

}

var setConfigDesc = "set current context config."

type cfgSetCommand struct {
	baseCommand
}

// Init initializes command.
func (cc *cfgSetCommand) Init(c *Cli) {
	cc.cli = c

	cc.cmd = &cobra.Command{
		Use:   "set cfgName",
		Short: "set current config",
		Long:  setConfigDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runSetConfig(args)
			}
		},
	}

	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *cfgSetCommand) addFlags() {
	// TODO: add flags here
}

func (cc *cfgSetCommand) runSetConfig(args []string) error {
	return myssh.SetContext(args[0], &cc.cli.Env)
}
