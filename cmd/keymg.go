package main

import (
	"fmt"
	"github.com/cnlubo/myssh/myssh"
	"github.com/cnlubo/myssh/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

var keyManageDesc = "Manage multiple SSH keys."

type KeyCommand struct {
	baseCommand
}

// Init initialize command.
func (cc *KeyCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "km",
		Aliases: []string{"mkm"},
		Short:   "manage ssh keys (alias: mkm)",
		Long:    keyManageDesc,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if found := utils.PathExist(cc.cli.Env.SKMPath); !found {
				err := os.Mkdir(cc.cli.Env.SKMPath, 0755)
				if err != nil {
					return errors.Wrap(err, "Create SSH keyStore dir failed")
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.cmd.Help()
			}
		},
	}

	c.AddCommand(cc, &keyStoreInitCmd{})
	c.AddCommand(cc, &keyListCmd{})
	c.AddCommand(cc, &keyAddCmd{})
	c.AddCommand(cc, &keyDeleteCmd{})
	c.AddCommand(cc, &keySetCmd{})
	c.AddCommand(cc, &keyDisplayCmd{})
	c.AddCommand(cc, &keyRenameCmd{})
	c.AddCommand(cc, &keyCopyCmd{})
}

var keyStoreInitDesc = "Initialize SSH keys store for the first time usage."

type keyStoreInitCmd struct {
	baseCommand
}

func (cc *keyStoreInitCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize SSH keys store",
		Long:  keyStoreInitDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runKeyStoreInit()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keyStoreInitCmd) addFlags() {
	// TODO: add flags here
}

func (cc *keyStoreInitCmd) runKeyStoreInit() error {
	return myssh.KeyStoreInit(&cc.cli.Option.Env)
}

var keyListDesc = " List all available SSH keys."

type keyListCmd struct {
	baseCommand
}

func (cc *keyListCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available SSH keys (alias:ls)",
		Long:    keyListDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runSSHKeyList()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keyListCmd) addFlags() {
	// TODO: add flags here
}

var keyHeader = []string{
	"AliasName",
	" KeyType ",
	" KeyDesc ",
}

func (cc *keyListCmd) runSSHKeyList() error {

	keyMap, _ := myssh.LoadSSHKeys(&cc.cli.Option.Env)
	if len(keyMap) == 0 {
		return errors.New("No SSH key found")
	}
	utils.PrintN(utils.Info, fmt.Sprintf("Found %d SSH key(s)!\r\n", len(keyMap)))
	fmt.Println()
	cc.cli.PrintTable(displaySSHKey(keyMap))
	return nil
}

func displaySSHKey(keyMap map[string]*myssh.SSHKey) (tableHeader []string, hColors []tablewriter.Colors, colColors []tablewriter.Colors, data [][]string) {

	var keys []string
	for k := range keyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var rowData [][]string
	for _, k := range keys {
		key := keyMap[k]
		var row []string

		if key.IsDefault {
			row = append(row, utils.CheckSymbol+k, key.Type.Name, key.KeyDesc)
		} else {
			row = append(row, k, key.Type.Name, key.KeyDesc)
		}
		rowData = append(rowData, row)
	}
	var headerColors, columnColors []tablewriter.Colors
	for i := 1; i <= len(keyHeader); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgHiMagentaColor, tablewriter.Bold})
		columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgHiMagentaColor})
	}
	return keyHeader, headerColors, columnColors, rowData
}

var keyAddDesc = "Add one new SSHKey."

type keyAddCmd struct {
	baseCommand
}

// Init initialize command.
func (cc *keyAddCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "add",
		Short: "add one SSHKey",
		Long:  keyAddDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runCreateSSHKey()
			}
		},
	}
	cc.addFlags()
}

func (cc *keyAddCmd) addFlags() {
	// TODO
}

func (cc *keyAddCmd) runCreateSSHKey() error {

	keys, err := myssh.CreateSSHKey(&cc.cli.Option.Env)
	if err != nil {
		return err
	}
	fmt.Println()
	cc.cli.PrintTable(displaySSHKey(keys))
	return nil
}

var keyDeleteDesc = "delete specific SSHKey by alias name."

type keyDeleteCmd struct {
	baseCommand
}

func (cc *keyDeleteCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "delete alias ",
		Aliases: []string{"del"},
		Short:   "delete SSH key (alias:del)",
		Long:    keyDeleteDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runDeleteSSHKey(args[0])
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keyDeleteCmd) addFlags() {
	// TODO: add flags here
}

func (cc *keyDeleteCmd) runDeleteSSHKey(aliasName string) error {

	return myssh.DeleteSSHKey(aliasName, &cc.cli.Option.Env)
}

var keySetDesc = "Set specific SSH key as default by its alias name."

type keySetCmd struct {
	baseCommand
}

func (cc *keySetCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "use",
		Short: "Set specific SSH key as default",
		Long:  keySetDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runSetSSHKey(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keySetCmd) addFlags() {
	// TODO: add flags here
}

func (cc *keySetCmd) runSetSSHKey(args []string) error {

	var aliasName string
	if len(args) == 0 {
		aliasName = ""
	} else {
		aliasName = args[0]
	}

	keys, err := myssh.SetSSHKey(aliasName, &cc.cli.Option.Env)
	if err != nil {
		return err
	}
	fmt.Println()
	cc.cli.PrintTable(displaySSHKey(keys))
	return nil
}

var DisplayKeyDesc = "Display the current SSH public key or specific SSH public key by alias."

type keyDisplayCmd struct {
	baseCommand
}

func (cc *keyDisplayCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "display alias",
		Aliases: []string{"dp"},
		Short:   "Display SSHKey (alias:dp)",
		Long:    DisplayKeyDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runDisplaySSHKey(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keyDisplayCmd) addFlags() {
	// TODO: add flags here
}

func (cc *keyDisplayCmd) runDisplaySSHKey(args []string) error {

	var aliasName string
	if len(args) == 0 {
		aliasName = ""
	} else {
		aliasName = args[0]
	}
	return myssh.DisplaySSHKey(aliasName, &cc.cli.Option.Env)
}

var keyRenameDesc = "Rename SSH key aliasName."

type keyRenameCmd struct {
	baseCommand
}

func (cc *keyRenameCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "rename alias newalias",
		Aliases: []string{"rn"},
		Short:   "Rename SSH key aliasName (alias:rn)",
		Long:    keyRenameDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 2, 2); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runKeyRename(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keyRenameCmd) addFlags() {
	// TODO: add flags here
}

func (cc *keyRenameCmd) runKeyRename(args []string) error {

	return myssh.RenameSSHKey(args[0], args[1], &cc.cli.Option.Env)
}

var keyCopyDesc = "Copy SSH public key to a remote host."

type keyCopyCmd struct {
	baseCommand
	identityfile string
}

func (cc *keyCopyCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "copy host",
		Aliases: []string{"cp"},
		Short:   "Copy SSH public key to a remote host (alias:cp)",
		Long:    keyCopyDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "arguments input failed")
			} else {
				return cc.runSSHKeyCopy(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *keyCopyCmd) addFlags() {
	// TODO: add flags here
}
func (cc *keyCopyCmd) runSSHKeyCopy(args []string) error {
	// parse connect string, format is [user@]host[:port]
	u, h, p := utils.ParseConnect(args[0])
	if len(u) == 0 {
		return errors.New(fmt.Sprintf("bad host [%s] user not exists!!!! ", args[0]))
	}
	if len(p) == 0 {
		p = "22"
	}
	connectStr := u + "@" + h
	return myssh.CopySSHKey(connectStr, p, &cc.cli.Option.Env)
}
