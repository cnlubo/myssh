package main

import (
	"fmt"
	"github.com/cnlubo/myssh/myssh"
	"github.com/cnlubo/myssh/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sort"
)

var clusterDesc = "manage clusters."

type ClusterCommand struct {
	baseCommand
}

var tableHeader = []string{
	"Host List",
}

// Init initialize command.
func (cc *ClusterCommand) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "clusters",
		Aliases: []string{"mclusters"},
		Short:   "manage clusters (alias: mclusters)",
		Long:    clusterDesc,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//
		// 	return myssh.ClusterInit(&cc.cli.Env)
		// },
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

	c.AddCommand(cc, &parseCmd{})
	c.AddCommand(cc, &clusterAddCmd{})
	c.AddCommand(cc, &clusterListCmd{})
	c.AddCommand(cc, &clusterDelCmd{})
	c.AddCommand(cc, &clusterBatchCmd{})
	c.AddCommand(cc, &clusterKeyCopyCmd{})
	c.AddCommand(cc, &clusterCopyCmd{})

}

var parseDesc = "parse host pattern to host list."

type parseCmd struct {
	baseCommand
	expand bool
}

// Init initializes command.
func (cc *parseCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "parse HOST_PATTERN [flag]",
		Aliases: []string{"pa"},
		Short:   "parse hostPattern to host list (alias:pa)",
		Long:    parseDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runParseHosts(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *parseCmd) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.BoolVarP(&cc.expand, "expand", "e", false, "Expand the host list output to multiple lines")
}

func (cc *parseCmd) runParseHosts(args []string) error {

	hosts, err := myssh.ParseExpr(args[0])
	if err != nil {
		return err
	}

	sort.Sort(sort.StringSlice(hosts))

	if ok := cc.expand; ok {

		var rowData [][]string

		for _, v := range hosts {
			rowData = append(rowData, []string{v})
		}

		var headerColors, columnColors []tablewriter.Colors
		for i := 1; i <= len(tableHeader); i++ {
			headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold})
			columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgHiWhiteColor})

		}

		fmt.Println()

		cc.cli.PrintTable(tableHeader, headerColors, columnColors, rowData)

	} else {

		fmt.Println()

		cc.cli.PrintTable(nil, nil, nil, [][]string{hosts})

	}

	return nil

}

var clusterAddDesc = "Add one cluster."

type clusterAddCmd struct {
	baseCommand
}

// Init initializes command.
func (cc *clusterAddCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:   "add",
		Short: "Add one cluster",
		Long:  clusterAddDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runAddCluster()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *clusterAddCmd) addFlags() {
	// TODO add flags here
}

func (cc *clusterAddCmd) runAddCluster() error {

	clusters, err := myssh.AddCluster()
	if err != nil {
		return err
	}
	cc.cli.PrintTable(displayClusters(clusters))
	return nil

}

var clusterListDesc = "List all cluster."

type clusterListCmd struct {
	baseCommand
}

// Init initializes command.
func (cc *clusterListCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all cluster (alias:ls)",
		Long:    clusterListDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runListClusters()
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *clusterListCmd) addFlags() {
	// TODO add flags here
}

func (cc *clusterListCmd) runListClusters() error {

	clusters := myssh.ClustersCfg.Clusters
	if clusters.Len() == 0 {
		return errors.New("No cluster found")
	}
	fmt.Println()
	cc.cli.PrintTable(displayClusters(clusters))
	return nil
}

var clusterDeleteDesc = "delete cluster."

type clusterDelCmd struct {
	baseCommand
}

// Init initializes command.
func (cc *clusterDelCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "delete cluster...",
		Aliases: []string{"del"},
		Short:   "delete cluster (alias:del)",
		Long:    clusterDeleteDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, -1); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runDeleteClusters(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *clusterDelCmd) addFlags() {
	// TODO add flags here
}

func (cc *clusterDelCmd) runDeleteClusters(args []string) error {

	_, err := myssh.DeleteClusters(args)
	if err != nil {
		return err
	}

	return nil
}

var clusterKeyCopyDesc = "Copy public key to cluster."

type clusterKeyCopyCmd struct {
	baseCommand
	sshPort      int
	sshUser      string
	identityfile string
}

// Init initialize command.
func (cc *clusterKeyCopyCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "keycopy hostPatterns [flags]",
		Aliases: []string{"kcp"},
		Short:   "Copy public key to cluster (alias: kcp)",
		Long:    clusterKeyCopyDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runClusterKeyCopy(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *clusterKeyCopyCmd) addFlags() {

	flagSet := cc.cmd.Flags()

	flagSet.IntVarP(&cc.sshPort, "port", "p", 0, "Port for the remote SSH service")
	flagSet.StringVarP(&cc.sshUser, "user", "u", "", "User account for SSH login")
	flagSet.StringVarP(&cc.identityfile, "identityfile", "i", "", "identity file (private key) for public key authentication.")
}

func (cc *clusterKeyCopyCmd) runClusterKeyCopy(args []string) error {

	return myssh.ClusterKeyCopy(args[0], cc.sshPort, cc.sshUser, cc.identityfile)
}

var clusterBatchDesc = "Batch exec command for cluster."

type loginOption struct {
	Login myssh.ServerConfig
}

type clusterBatchCmd struct {
	baseCommand
	Prompt bool
	loginOption
}

// Init initialize command.
func (cc *clusterBatchCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "batch hostPatterns command ... [flags]",
		Aliases: []string{"bt"},
		Short:   "batch exec command (alias: bt)",
		Long:    clusterBatchDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 2, -1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runBatch(args)
			}
		},
		// Example: batchExample(),
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *clusterBatchCmd) addFlags() {

	flagSet := cc.cmd.Flags()
	flagSet.IntVarP(&cc.Login.Port, "port", "p", 22, "Port for the remote SSH service")
	flagSet.StringVarP(&cc.Login.User, "user", "u", "", "User account for SSH login")
	flagSet.StringVarP(&cc.Login.PrivateKey, "identityfile", "i", "", "identity file (private key) for public key authentication.")
	// flagSet.StringVarP(&cc.Login.PrivateKeyPassword, "identitypass", "P", "", "identity password (private key) for public key authentication.")
	flagSet.BoolVarP(&cc.Prompt, "password", "P", false, "Prompt for password")
}

func (cc *clusterBatchCmd) runBatch(args []string) error {
	return myssh.ClusterBatchCmds(args[0], cc.Prompt, &cc.Login, args[1:]...)
}

// execExample shows examples in exec command, and is used in auto-generated cli docs.
// func batchExample() string {
// 	return `$ myssh cluster batch hostPattern cmd...`
// }

func displayClusters(clusters myssh.Clusters) (tableHeader []string, hColors []tablewriter.Colors, colColors []tablewriter.Colors, data [][]string) {

	var clusterTableHeader = []string{
		" ClusterName  ",
		" HostPattern ",
	}

	var rowData [][]string
	for _, v := range clusters {
		rowData = append(rowData, []string{v.Name, v.HostPattern})
	}

	var headerColors, columnColors []tablewriter.Colors
	for i := 1; i <= len(clusterTableHeader); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold})
		columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgGreenColor})

	}
	return clusterTableHeader, headerColors, columnColors, rowData

}

var ClusterCopyDesc = "Copy files or Directory to cluster."

type clusterCopyCmd struct {
	baseCommand
	Prompt bool
	loginOption
}

// Init initializes command.
func (cc *clusterCopyCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "copy hostPatterns source target",
		Aliases: []string{"mcp"},
		Short:   "copy files or Directory to cluster hosts (alias:mcp)",
		Long:    ClusterCopyDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), 3, 3); err != nil {
				myssh.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				// return cc.runParseHosts(args)
				return nil
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *clusterCopyCmd) addFlags() {
	// flagSet.BoolVarP(&cc.singleCPServer, "single", "s", false, "single server")
	flagSet := cc.cmd.Flags()
	flagSet.IntVarP(&cc.Login.Port, "port", "p", 22, "Port for the remote SSH service")
	flagSet.StringVarP(&cc.Login.User, "user", "u", "", "User account for SSH login")
	flagSet.StringVarP(&cc.Login.PrivateKey, "identityfile", "i", "", "identity file (private key) for public key authentication.")
	// flagSet.StringVarP(&cc.Login.PrivateKeyPassword, "identitypass", "P", "", "identity password (private key) for public key authentication.")
	flagSet.BoolVarP(&cc.Prompt, "password", "P", false, "Prompt for password")
}

func (cc *clusterCopyCmd) runCopy(args []string) error {
	// return myssh.ClusterCopy(args, cc.Prompt, &cc.Login)
	return nil

}
