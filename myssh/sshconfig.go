package myssh

import (
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"github.com/cnlubo/promptx"
	"github.com/cnlubo/ssh_config"
	"github.com/goinggo/mapstructure"
	"github.com/gookit/color"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// HostConfig struct include alias, connect string and other config
type HostConfig struct {
	// Alias alias
	Alias string
	// Path found in which file
	Path string
	// PathMap key is file path, value is the alias's hosts
	PathMap map[string][]*ssh_config.Host
	// OwnConfig own config
	OwnConfig map[string]string
	// ImplicitConfig implicit config
	ImplicitConfig map[string]string
}

// List ssh alias, filter by optional keyword
func ListAlias(p string, lo ListOption) ([]*HostConfig, error) {

	configMap, aliasMap, err := parseConfig(p)
	// fmt.Println(aliasMap["test"].PathMap)
	if err != nil {
		return nil, err
	}
	var result []*HostConfig
	for _, host := range aliasMap {
		values := []string{host.Alias}
		for _, v := range host.OwnConfig {
			values = append(values, v)
		}
		if len(lo.Keywords) > 0 && !utils.Query(values, lo.Keywords, lo.IgnoreCase) {
			continue
		}
		result = append(result, host)
	}

	// Format
	for fp, cfg := range configMap {
		if len(cfg.Hosts) > 0 {
			if err := writeConfig(fp, cfg); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func DeleteAliases(p string, aliases ...string) ([]*HostConfig, error) {

	configMap, aliasMap, err := parseConfig(p)
	if err != nil {
		return nil, err
	}
	// check alias exists
	if err := checkAlias(aliasMap, true, aliases...); err != nil {
		return nil, err
	}

	var deleteHosts []*HostConfig
	var confirm bool
	clName := strings.Join(aliases, ",")
	message := "(Please confirm to delete aliases [" + clName + "]"
	c := promptx.NewDefaultConfirm(message, true)
	confirm, err = c.Run()
	if err != nil {
		return nil, nil
	}
	if confirm {
		for _, alias := range aliases {
			deleteHost := aliasMap[alias]
			deleteHosts = append(deleteHosts, deleteHost)
			for fp, hosts := range deleteHost.PathMap {
				for _, host := range hosts {
					if len(host.Patterns) == 1 {
						deleteHostFromConfig(configMap[fp], host)
					} else {
						var patterns []*ssh_config.Pattern
						for _, pattern := range host.Patterns {
							if pattern.String() != alias {
								patterns = append(patterns, pattern)
							}
						}
						host.Patterns = patterns
					}
				}
				if err := writeConfig(fp, configMap[fp]); err != nil {
					return nil, err
				}
			}
		}

		return deleteHosts, nil
	} else {
		return nil, nil
	}

}

type addOption struct {
	// Path add path
	Path string
	// Alias alias
	Alias string
	// Connect connection string
	Connect string
	// Config other config
	Config map[string]string
}

// Add ssh host config to ssh config file
func addAlias(p string, ao *addOption) (*HostConfig, error) {

	if ao.Path == "" {
		ao.Path = p
	}

	configMap, aliasMap, err := parseConfig(p)
	if err != nil {
		return nil, err
	}
	// check alias exists
	if err := checkAlias(aliasMap, false, ao.Alias); err != nil {
		return nil, err
	}

	cfg, ok := configMap[ao.Path]
	if !ok {
		cfg, err = readCfgFile(ao.Path)
		if err != nil {
			return nil, err
		}
	}

	// Parse connect string
	user, hostname, port := utils.ParseConnect(ao.Connect)
	if user != "" {
		ao.Config["user"] = user
	}
	if hostname != "" {
		ao.Config["hostname"] = hostname
	}
	if port != "" {
		ao.Config["port"] = port
	}

	var nodes []ssh_config.Node
	for k, v := range ao.Config {
		nodes = append(nodes, NewKV(strings.ToLower(k), v))
	}

	pattern, err := ssh_config.NewPattern(ao.Alias)
	if err != nil {
		return nil, err
	}

	cfg.Hosts = append(cfg.Hosts, &ssh_config.Host{
		Patterns: []*ssh_config.Pattern{pattern},
		Nodes:    nodes,
	})

	_ = writeConfig(ao.Path, cfg)

	_, aliasMap, err = parseConfig(p)
	if err != nil {
		return nil, err
	}
	return aliasMap[ao.Alias], nil
}

func AliasKeyCopy(hostCfgPath string, aliases []string) error {

	if len(aliases) == 0 {
		return errors.New("alias must been input")
	}
	servers, err := loadAlias(hostCfgPath, false, ListOption{
		Keywords:   aliases,
		IgnoreCase: true,
	})
	if err != nil {
		return err
	}

	for _, server := range servers {

		if len(utils.CompressStr(server.Address)) > 0 {

			var identityfile, connectStr, port, sshUser string

			if len(utils.CompressStr(server.PrivateKey)) > 0 {
				identityfile = server.PrivateKey
				home, _ := homedir.Dir()
				identityfile = utils.ParseRelPath(identityfile, home)
				// check identityfile
				_, err := privateKeyFile(identityfile, " ")
				if err != nil {
					utils.PrintErr(errors.Wrap(err, fmt.Sprintf("alias (%s) have bad identityfile", server.Name)))
					continue
				}
			} else {
				utils.PrintN(utils.Err, fmt.Sprintf("alias (%s) have empty identityfile", server.Name))
				continue
			}

			if server.Port == 0 {
				port = "22"
			} else {
				port = strconv.Itoa(server.Port)
			}

			if len(utils.CompressStr(server.User)) == 0 {
				sshUser = utils.GetUsername()
			} else {
				sshUser = server.User
			}

			connectStr = sshUser + "@" + utils.CompressStr(server.Address)
			utils.PrintN(utils.Info, fmt.Sprintf("copy current SSH key to [%s]\n", server.Name))
			result := CopyKey(connectStr, port, identityfile)
			if result {
				utils.PrintN(utils.Info, fmt.Sprintf("Current SSH key have been copied to [%s]\n", server.Name))
			}
		}
	}

	return nil
}

// UpdateOption options for Update
type updateOption struct {
	// Alias alias
	Alias string
	// NewAlias new alias
	NewAlias string
	// Connect connection string
	Connect string
	// Config other config
	Config map[string]string
}

func UpdateHostCfg(hostCfgPath string, aliasName string) ([]*HostConfig, error) {

	if len(aliasName) == 0 {
		return nil, errors.New("alias must been input")
	}
	var hostCfg *HostConfig
	if hostCfg = findAliasByName(hostCfgPath, aliasName); hostCfg == nil {
		return nil, errors.New(fmt.Sprintf("ssh host alias (%s) not exists", aliasName))
	}

	uo := &updateOption{
		Alias: aliasName,
	}

	p := promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		} else {
			if err := utils.CheckConnect(string(line)); err != nil {
				return err
			}
		}
		return nil

	}, "ConnectString:")

	p.Default = hostCfg.ConnectionStr()

	connectString, err := p.Run()
	if err != nil {
		return nil, err
	}

	// identityfile
	p = promptx.NewDefaultPrompt(func(line []rune) error {
		// allow empty
		return nil
	}, "Identityfile:")

	if value, ok := hostCfg.OwnConfig["identityfile"]; ok {
		home, _ := homedir.Dir()
		p.Default = utils.ParseRelPath(value, home)
	}

	identityfile, err := p.Run()
	if err != nil {
		return nil, err
	}

	// new aliasName

	p = promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			// allow empty
			return nil
		} else {
			// new aliasName should not exists
			if s := findAliasByName(hostCfgPath, string(line)); s != nil {
				return errors.New(fmt.Sprintf("NewAlias (%s) exists", string(line)))
			}
		}
		return nil
	}, "NewAliasName:")

	newAliasName, err := p.Run()
	if err != nil {
		return nil, err
	}

	// set updateOption
	uo.Connect = connectString

	uo.Config = make(map[string]string)
	if len(utils.CompressStr(identityfile)) != 0 {
		uo.Config["identityfile"] = identityfile
	}
	if len(utils.CompressStr(newAliasName)) != 0 {
		uo.NewAlias = newAliasName
	}
	if !uo.valid() {
		return nil, errors.New("the update option is invalid")
	}

	host, err := updateCfg(hostCfgPath, uo)
	if err != nil {
		return nil, errors.Wrap(err, "update alias failed")
	}
	utils.PrintN(utils.Info, "update successfully\n")
	if host != nil {
		return []*HostConfig{host}, nil
	}
	return nil, nil
}

func updateCfg(p string, uo *updateOption) (*HostConfig, error) {

	configMap, aliasMap, err := parseConfig(p)
	if err != nil {
		return nil, err
	}

	if err := checkAlias(aliasMap, true, uo.Alias); err != nil {
		return nil, errors.Wrapf(err, "alias %s not found", uo.Alias)
	}
	updateHost := aliasMap[uo.Alias]

	if uo.NewAlias != "" {
		// new alias should not exist
		if err := checkAlias(aliasMap, false, uo.NewAlias); err != nil {
			return nil, err
		}
	} else {
		uo.NewAlias = uo.Alias
	}

	if uo.Connect != "" {
		// Parse connect string
		user, hostname, port := utils.ParseConnect(uo.Connect)
		if user != "" {
			uo.Config["user"] = user
		}
		if hostname != "" {
			uo.Config["hostname"] = hostname
		}
		if port != "" {
			uo.Config["port"] = port
		}
	}
	// update configs
	for k, v := range uo.Config {
		if v == "" {
			delete(updateHost.OwnConfig, k)
		} else {
			updateHost.OwnConfig[k] = v
		}
	}

	// update host
	for fp, hosts := range updateHost.PathMap {
		for i, host := range hosts {
			if fp == updateHost.Path {
				pattern, _ := ssh_config.NewPattern(uo.NewAlias)
				newHost := &ssh_config.Host{
					Patterns: []*ssh_config.Pattern{pattern},
				}
				for k, v := range updateHost.OwnConfig {
					newHost.Nodes = append(newHost.Nodes, NewKV(k, v))
				}
				if len(host.Patterns) == 1 {
					if i == 0 {
						*host = *newHost
						// for implicit "*"
						find := false
						for _, h := range configMap[fp].Hosts {
							if host == h {
								find = true
								break
							}
						}
						if !find {
							newHost.Nodes = []ssh_config.Node{}
							for k, v := range uo.Config {
								newHost.Nodes = append(newHost.Nodes, NewKV(k, v))
							}
							configMap[fp].Hosts = append(configMap[fp].Hosts, newHost)
						}
					} else {
						deleteHostFromConfig(configMap[fp], host)
					}
				} else {
					if i == 0 {
						configMap[fp].Hosts = append(configMap[fp].Hosts, newHost)
					}
					var patterns []*ssh_config.Pattern
					for _, pattern := range host.Patterns {
						if pattern.String() != uo.NewAlias {
							patterns = append(patterns, pattern)
						}
					}
					host.Patterns = patterns
				}
			} else {
				if len(host.Patterns) == 1 {
					deleteHostFromConfig(configMap[fp], host)
				} else {
					var patterns []*ssh_config.Pattern
					for _, pattern := range host.Patterns {
						if pattern.String() != uo.NewAlias {
							patterns = append(patterns, pattern)
						}
					}
					host.Patterns = patterns
				}
			}
			if err := writeConfig(fp, configMap[fp]); err != nil {
				return nil, err
			}
		}
	}
	_, aliasMap, err = parseConfig(p)
	if err != nil {
		return nil, err
	}
	return aliasMap[uo.NewAlias], nil
}

func AddHostAlias(hostCfgPath string) ([]*HostConfig, error) {

	sshConfigFile := hostCfgPath
	// alias
	p := promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		} else if len(line) > 12 {
			return inputTooLongErr
		}

		if s := findAliasByName(sshConfigFile, string(line)); s != nil {
			return errors.New(fmt.Sprintf("Alias (%s) exists", string(line)))
		}
		return nil

	}, "Alias:")

	alias, err := p.Run()
	if err != nil {
		return nil, err
	}

	// connectString
	p = promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		} else {
			if err := utils.CheckConnect(string(line)); err != nil {
				return err
			}
		}
		return nil

	}, "ConnectString:")

	connectString, err := p.Run()
	if err != nil {
		return nil, err
	}

	p = promptx.NewDefaultPrompt(func(line []rune) error {
		// allow empty
		return nil

	}, "Identityfile:")

	identityfile, err := p.Run()
	if err != nil {
		return nil, err
	}

	ao := &addOption{
		Alias:   alias,
		Connect: connectString,
		Path:    sshConfigFile,
	}

	ao.Config = make(map[string]string)
	if identityfile != "" {
		ao.Config["identityfile"] = identityfile
	}
	host, err := addAlias(sshConfigFile, ao)
	if err != nil {
		return nil, errors.Wrap(err, "add alias failed")
	}
	utils.PrintN(utils.Info, "added alias successfully\n")
	if host != nil {
		return []*HostConfig{host}, nil
	}
	return nil, nil
}

func BatchExecAlias(sshConfigPath string, promptPass bool, alias string, cmd ...string) error {

	reg := regexp.MustCompile(`[,]`)
	aliasList := reg.Split(alias, -1)
	servers, err := loadAlias(sshConfigPath, promptPass, ListOption{
		Keywords:   aliasList,
		IgnoreCase: true})
	if err != nil {
		return err
	}
	exCmd := strings.Join(cmd, "&&")
	if promptPass {
		utils.Clear()
	}
	fmt.Println()

	err = batchExec(servers, exCmd)
	if err != nil {
		return err
	}
	return nil
}

func AliasLogin(aliasName string, promptPass bool, sshConfigFile string) error {

	hostCfg := findAliasByName(sshConfigFile, utils.CompressStr(aliasName))

	var privateKey, loginPass, method string
	var server *ServerConfig
	if hostCfg != nil {
		if value, ok := hostCfg.OwnConfig["identityfile"]; ok {
			if result := path.IsAbs(value); !result {
				home, _ := homedir.Dir()
				privateKey = utils.ParseRelPath(value, home)
			} else {
				privateKey = value
			}
		} else {
			privateKey = ""
		}
		if len(privateKey) == 0 || promptPass {

			utils.PrintN(utils.Warn, "Login into "+hostCfg.OwnConfig["user"]+"@"+hostCfg.OwnConfig["hostname"])
			p := promptx.NewDefaultPrompt(func(line []rune) error {
				if strings.TrimSpace(string(line)) == "" {
					return errors.New("password is empty")
				}
				return nil
			}, "Password:")
			p.Mask = MaskPrompt
			loginPass, _ = p.Run()
			method = "password"

		} else {
			loginPass = ""
			method = "key"
		}
		var pp int
		if len(utils.CompressStr(hostCfg.OwnConfig["port"])) == 0 {
			pp = 22
		} else {
			pp, _ = strconv.Atoi(hostCfg.OwnConfig["port"])
		}
		server = &ServerConfig{
			Name:               hostCfg.Alias,
			User:               hostCfg.OwnConfig["user"],
			Address:            hostCfg.OwnConfig["hostname"],
			Method:             method,
			Port:               pp,
			PrivateKey:         privateKey,
			PrivateKeyPassword: "",
			Password:           loginPass,
		}
	} else {
		return errors.New(fmt.Sprintf("Server (%s) not found", aliasName))
	}

	return serverLogin(server)
}

func AliasInteractiveLogin(sshConfigFile string, promptPass bool) error {

	servers, err := loadAlias(sshConfigFile, false, ListOption{})
	if err != nil {
		return err
	}
	if servers == nil {
		// utils.PrintN(utils.Info, "not found ssh servers\n")
		return errors.New("not found ssh servers\n")
	}
	cfg := &promptx.SelectConfig{
		ActiveTpl:    `»  {{ .Name | cyan }}: {{ .User | cyan }}{{ "@" | cyan }}{{ .Address | cyan }}`,
		InactiveTpl:  `  {{ .Name | white }}: {{ .User | white }}{{ "@" | white }}{{ .Address | white }}`,
		SelectPrompt: "Login Server",
		SelectedTpl:  `{{ "» " | green }}{{ .Name | green }}: {{ .User | green }}{{ "@" | green }}{{ .Address | green }}`,
		DisPlaySize:  9,
		DetailsTpl: `
  --------- Login Server -----------------------
	{{ "Name:" | faint }} {{ .Name | faint }}
	{{ "User:" | faint }} {{ .User | faint }}
  {{ "PrivateKey:" | faint }} {{ .PrivateKey | faint }}
     {{ "Address:" | faint }} {{ .Address | faint }}{{ ":" | faint }}{{ .Port | faint }}`,
	}

	s := &promptx.Select{
		Items:  servers,
		Config: cfg,
	}
	idx := s.Run()
	loginServer := servers[idx]
	var loginPass string
	if len(loginServer.PrivateKey) == 0 || promptPass {
		utils.PrintN(utils.Warn, "Login into "+loginServer.User+"@"+loginServer.Address)
		p := promptx.NewDefaultPrompt(func(line []rune) error {
			if strings.TrimSpace(string(line)) == "" {
				return errors.New("password is empty")
			}
			return nil
		}, "LoginPass:")
		p.Mask = MaskPrompt
		loginPass, err = p.Run()
		if err != nil {
			return err
		}
		loginServer.Password = loginPass
		loginServer.Method = "password"
	} else {
		loginServer.Method = "key"
	}
	return serverLogin(loginServer)
}

// Display Whether to display connection string
func (hc *HostConfig) Display() bool {
	hostname := hc.OwnConfig["hostname"]
	if hostname == "" {
		hostname = hc.ImplicitConfig["hostname"]
	}

	return hostname != ""
}

// ConnectionStr return the connection string
func (hc *HostConfig) ConnectionStr() string {
	if !hc.Display() {
		return ""
	}

	var (
		user, hostname, port string
		ok                   bool
	)

	if user, ok = hc.OwnConfig["user"]; !ok {
		user = hc.ImplicitConfig["user"]
		delete(hc.ImplicitConfig, "user")
	} else {
		delete(hc.OwnConfig, "user")
	}

	if hostname, ok = hc.OwnConfig["hostname"]; !ok {
		delete(hc.ImplicitConfig, "hostname")
		hostname = hc.ImplicitConfig["hostname"]
	} else {
		delete(hc.OwnConfig, "hostname")
	}

	if port, ok = hc.OwnConfig["port"]; !ok {
		port = hc.ImplicitConfig["port"]
		delete(hc.ImplicitConfig, "port")
	} else {
		delete(hc.OwnConfig, "port")
	}
	return color.New(color.FgGreen, color.OpBold).Sprintf("%s%s%s%s%s", user, "@", hostname, ":", port)
}

// NewHostConfig new HostConfig
func newHostConfig(alias, path string, host *ssh_config.Host) *HostConfig {
	return &HostConfig{
		Alias:          alias,
		Path:           path,
		PathMap:        map[string][]*ssh_config.Host{path: {host}},
		OwnConfig:      map[string]string{},
		ImplicitConfig: map[string]string{},
	}
}

// Valid whether the option is valid
func (uo *updateOption) valid() bool {
	return uo.NewAlias != "" || uo.Connect != "" || len(uo.Config) > 0
}

func deleteHostFromConfig(config *ssh_config.Config, host *ssh_config.Host) {
	var hs []*ssh_config.Host
	for _, h := range config.Hosts {
		if h == host {
			continue
		}
		hs = append(hs, h)
	}
	config.Hosts = hs
}

// find alias by name
func findAliasByName(p string, name string) *HostConfig {

	_, aliasMap, err := parseConfig(p)
	if err != nil {
		return nil
	}
	for _, host := range aliasMap {
		if host.Alias == name {
			return host
		}
	}
	return nil
}

func checkAlias(aliasMap map[string]*HostConfig, expectExist bool, aliases ...string) error {

	for _, alias := range aliases {
		ok := aliasMap[alias] != nil
		if !ok && expectExist {
			return fmt.Errorf("alias [%s] not exists", alias)
		} else if ok && !expectExist {
			return fmt.Errorf("alias[%s] already exists", alias)
		}
	}
	return nil
}

// get alias return servers
func loadAlias(sshConfigPath string, promptPass bool, lo ListOption) (Servers, error) {
	type srv struct {
		Alias        string
		HostName     string
		User         string
		Port         string
		Identityfile string
	}
	var srvs Servers
	var loginPass, authMethod string

	_, aliasMap, err := parseConfig(sshConfigPath)
	if err != nil {
		return nil, err
	}
	for _, host := range aliasMap {
		var result srv
		values := []string{host.Alias}

		for _, v := range host.OwnConfig {
			values = append(values, v)
		}
		if len(lo.Keywords) > 0 && !utils.Query(values, lo.Keywords, lo.IgnoreCase) {
			continue
		}

		if err := mapstructure.Decode(host.OwnConfig, &result); err != nil {
			return nil, err
		}
		if result.HostName != "" {
			result.Alias = host.Alias
			pp, _ := strconv.Atoi(result.Port)
			if ok := path.IsAbs(result.Identityfile); !ok {
				home, _ := homedir.Dir()
				result.Identityfile = strings.Replace(result.Identityfile, "~", home, 1)
			}
			if promptPass {
				utils.PrintN(utils.Warn, "["+result.Alias+"] -> "+result.User+"@"+result.HostName)

				p := promptx.NewDefaultPrompt(func(line []rune) error {
					if strings.TrimSpace(string(line)) == "" {
						return errors.New("password is empty")
					}
					return nil
				}, "LoginPass:")

				p.Mask = MaskPrompt
				loginPass, err = p.Run()
				if err != nil {
					return nil, err
				}
				authMethod = "password"

			} else {
				loginPass = ""
				authMethod = "key"
			}
			s := &ServerConfig{
				Name:               result.Alias,
				User:               result.User,
				Address:            result.HostName,
				Port:               pp,
				PrivateKey:         result.Identityfile,
				PrivateKeyPassword: "",
				Password:           loginPass,
				Method:             authMethod,
			}
			srvs = append(srvs, s)
		}
	}
	if srvs != nil {
		return srvs, nil
	} else {
		return nil, nil
	}
}

func writeConfig(p string, cfg *ssh_config.Config) error {
	return ioutil.WriteFile(p, []byte(cfg.String()), 0644)
}
func setImplicitConfig(aliasMap map[string]*HostConfig, hc *HostConfig) {
	for alias, host := range aliasMap {
		if alias == hc.Alias {
			continue
		}

		if len(hc.OwnConfig) == 0 {
			if match, err := path.Match(host.Alias, hc.Alias); err != nil || !match {
				continue
			}
			for k, v := range host.OwnConfig {
				if _, ok := hc.ImplicitConfig[k]; !ok {
					hc.ImplicitConfig[k] = v
				}
			}
			continue
		}
		if match, err := path.Match(hc.Alias, host.Alias); err != nil || !match {
			continue
		}
		for k, v := range hc.OwnConfig {
			if _, ok := host.OwnConfig[k]; ok {
				continue
			}
			if _, ok := host.ImplicitConfig[k]; !ok {
				host.ImplicitConfig[k] = v
			}
		}
	}
}

func setOwnConfig(aliasMap map[string]*HostConfig, hc *HostConfig, h *ssh_config.Host) {
	if host, ok := aliasMap[hc.Alias]; ok {
		if _, ok := host.PathMap[hc.Path]; !ok {
			host.PathMap[hc.Path] = []*ssh_config.Host{}
		}
		host.PathMap[hc.Path] = append(host.PathMap[hc.Path], h)
		for k, v := range hc.OwnConfig {
			if _, ok := host.OwnConfig[k]; !ok {
				host.OwnConfig[k] = v
			}
		}
	} else {
		aliasMap[hc.Alias] = hc
	}
}

func addHosts(aliasMap map[string]*HostConfig, fp string, hosts ...*ssh_config.Host) {

	for _, host := range hosts {
		// except implicit `*`
		if len(host.Nodes) == 0 {
			continue
		}
		for _, pattern := range host.Patterns {
			alias := pattern.String()
			hc := newHostConfig(alias, fp, host)
			setImplicitConfig(aliasMap, hc)

			for _, node := range host.Nodes {
				if kvNode, ok := node.(*ssh_config.KV); ok {
					kvNode.Key = strings.ToLower(kvNode.Key)
					if _, ok := hc.ImplicitConfig[kvNode.Key]; !ok {
						hc.OwnConfig[kvNode.Key] = kvNode.Value
					}
				}
			}

			setImplicitConfig(aliasMap, hc)
			setOwnConfig(aliasMap, hc, host)
		}
	}
}

func readCfgFile(p string) (*ssh_config.Config, error) {

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return ssh_config.Decode(f)
}

// ParseConfig parse configs from ssh config file, return config object and alias map
func parseConfig(p string) (map[string]*ssh_config.Config, map[string]*HostConfig, error) {

	cfg, err := readCfgFile(p)
	if err != nil {
		return nil, nil, err
	}
	aliasMap := map[string]*HostConfig{}
	configMap := map[string]*ssh_config.Config{p: cfg}

	for _, host := range cfg.Hosts {

		for _, node := range host.Nodes {

			// switch t := node.(type) {
			// case *ssh_config.Include:
			// 	for fp, config := range t.GetFiles() {
			// 		configMap[fp] = config
			// 		addHosts(aliasMap, fp, config.Hosts...)
			// 	}
			//
			// case *ssh_config.KV:
			// 	addHosts(aliasMap, p, host)
			// 	// case *ssh_config.Empty:
			// }
			if _, ok := node.(*ssh_config.KV); ok {
				addHosts(aliasMap, p, host)
				break
			} else {
				if t, ok := node.(*ssh_config.Include); ok {
					for fp, config := range t.GetFiles() {
						configMap[fp] = config
						addHosts(aliasMap, fp, config.Hosts...)
					}
				}
			}

		}

	}
	// fmt.Println(aliasMap["test1"].PathMap)
	// addHosts(aliasMap, p, &ssh_config.Host{
	// 	Patterns: []*ssh_config.Pattern{(&ssh_config.Pattern{}).SetStr("*")},
	// 	Nodes: []ssh_config.Node{
	// 		NewKV("user", utils.GetUsername()),
	// 		NewKV("port", "22"),
	// 	},
	// })
	// fmt.Println(len(aliasMap))
	// fmt.Println(aliasMap["cameochina"])
	// fmt.Println(aliasMap["test"].PathMap["/Users/ak47/.ssh/config"])
	return configMap, aliasMap, nil
}

func NewKV(key, value string) *ssh_config.KV {
	k := ssh_config.KV{}
	k.SetLeadingSpace(4)
	k.Key = key
	k.Value = value
	return &k
}
