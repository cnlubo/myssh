package myssh

import (
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"github.com/cnlubo/promptx"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	Main MainConfig
)

// set main config file path

func (cfg *MainConfig) SetConfigPath(configPath string) {
	cfg.configPath = configPath
}
func (cfg *MainConfig) GetConfigPath() string {
	return cfg.configPath
}

// write main config
func (cfg *MainConfig) write() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfg.configPath, out, 0644)
}

// write main config to yaml file
func (cfg *MainConfig) WriteTo(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.write()
}

// load main config
func (cfg *MainConfig) load() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	buf, err := ioutil.ReadFile(cfg.configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, cfg)
}

// load config from yaml file
func (cfg *MainConfig) LoadFrom(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.load()
}

// main config example
func MainConfigExample(env *Environment) *MainConfig {

	return &MainConfig{
		Basic: "default",
		Contexts: []Context{
			{
				Name:          "default",
				SSHConfig:     filepath.Join(env.StorePath, "contexts", "default", "sshconfig"),
				ClusterConfig: filepath.Join(env.StorePath, "contexts", "default", "default.yaml"),
			},
		},
		Current: "default",
	}
}

// delete context
func DeleteConfig(cfgs []string, env *Environment) (Contexts, error) {

	var deletesIdx []int
	var deleteCfgs Contexts
	var deleteNames []string
	for _, cfg := range cfgs {
		for i, s := range Main.Contexts {
			if strings.ToLower(s.Name) == strings.ToLower(cfg) {
				if Main.Current != s.Name {
					deletesIdx = append(deletesIdx, i)
					deleteCfgs = append(deleteCfgs, s)
					deleteNames = append(deleteNames, s.Name)
				} else {
					utils.PrintN(utils.Warn, fmt.Sprintf("[%s] is currently in use not allow delete !!!", s.Name))
				}
			}
		}
	}

	var confirm bool
	var err error
	cfName := strings.Join(deleteNames, ",")
	if len(deletesIdx) == 0 {
		return nil, errors.New("none config delete!!!")
	} else {
		prompt := promptx.NewDefaultConfirm("Please confirm to delete config ["+cfName+"]", true)
		confirm, err = prompt.Run()
		if err != nil {
			return nil, err
		}

	}

	if confirm {
		// sort and delete
		sort.Ints(deletesIdx)
		for i, del := range deletesIdx {
			Main.Contexts = append(Main.Contexts[:del-i], Main.Contexts[del-i+1:]...)
		}
		// save config
		sort.Sort(Main.Contexts)
		err := Main.write()
		if err != nil {
			return nil, errors.Wrap(err, "delete config failed")
		}
		// del config file
		for _, cf := range deleteCfgs {
			if utils.PathExist(filepath.Join(env.StorePath, "contexts", cf.Name)) {
				_ = os.RemoveAll(filepath.Join(env.StorePath, "contexts", cf.Name))
			}
		}
		utils.PrintN(utils.Info, fmt.Sprintf("delete config successfully\n"))
	} else {
		return nil, nil
	}
	return deleteCfgs, nil
}

func AddConfig(env *Environment) (Contexts, error) {

	var err error
	var ctxName string
	// CtxName
	prompt := promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return inputEmptyErr
		} else {
			if len(strings.TrimSpace(string(line))) > 12 {
				return inputTooLongErr
			}
		}
		_, ok := Main.Contexts.FindContextByName(string(line))
		if ok {
			return errors.New(fmt.Sprintf("context [%s] already existed", string(line)))
		}
		return nil

	}, "ContextName:")

	if ctxName, err = prompt.Run(); err != nil {
		return nil, err
	}

	cfgPath := filepath.Join(env.StorePath, "contexts", ctxName)
	if utils.PathExist(cfgPath) {
		_ = os.RemoveAll(cfgPath)
	}
	if err = os.MkdirAll(cfgPath, 0755); err != nil {
		return nil, errors.Wrap(err, "Create context dir error")
	}

	ctxClusterCfg := filepath.Join(cfgPath, ctxName+".yaml")
	if exists := utils.PathExist(ctxClusterCfg); !exists {

		utils.PrintN(utils.Info, fmt.Sprintf(" Create config file %s\n", ctxClusterCfg))

		if err := ClustersConfigExample(env).WriteTo(ctxClusterCfg); err != nil {
			return nil, errors.Wrap(err, "create configfile failed")
		}
	}
	ctxSSHCfg := filepath.Join(cfgPath, "sshconfig")
	if exists := utils.PathExist(ctxSSHCfg); !exists {

		utils.PrintN(utils.Info, fmt.Sprintf(" Create SSH config %s\n", ctxSSHCfg))
		// create ssh config file
		if err := ioutil.WriteFile(ctxSSHCfg, []byte("Include include/*"+"\n"), 0777); err != nil {
			return nil, errors.Wrap(err, "Create SSH config file error")
		}
		if err := os.MkdirAll(filepath.Join(cfgPath, "include"), 0755); err != nil {
			return nil, errors.Wrap(err, "Create context include dir error")
		}
	}

	// create config
	context := Context{
		Name:          ctxName,
		SSHConfig:     ctxSSHCfg,
		ClusterConfig: ctxClusterCfg,
	}
	// Save
	Main.Contexts = append(Main.Contexts, context)
	sort.Sort(Main.Contexts)
	err = Main.write()
	if err != nil {
		return nil, errors.Wrap(err, "add context failed")
	}
	utils.PrintN(utils.Info, fmt.Sprintf("Add Context[%s] success\n", ctxName))
	return Contexts{context}, nil

}

// find context by name
func (cs Contexts) FindContextByName(name string) (Context, bool) {
	for _, ctx := range cs {
		if name == ctx.Name {
			return ctx, true
		}
	}
	return Context{}, false
}

// set current context
func SetContext(cfgName string, env *Environment) error {

	if _, ok := Main.Contexts.FindContextByName(cfgName); !ok {
		return errors.New(fmt.Sprintf("config [%s] not exists", cfgName))
	}
	Main.Current = cfgName
	if err := Main.write(); err != nil {
		return errors.Wrap(err, "setup current config failed")
	}

	if err := CreateSSHlink(cfgName, env); err != nil {
		return errors.Wrap(err, "Create symlink error")
	}

	utils.PrintN(utils.Info, fmt.Sprintf("Now using Context:%s\n", cfgName))
	return nil
}

func CreateSSHlink(ctxName string, env *Environment) error {

	var (
		isLink    bool
		isLinkdir bool
		err       error
	)

	sshCfgFile := filepath.Join(env.SSHPath, "config")
	ctxSSHCfgFile := filepath.Join(env.StorePath, "contexts", ctxName, "sshconfig")

	if err, isLink = utils.IsSymLink(sshCfgFile); err != nil {
		// return err
		isLink = false
	}

	if utils.PathExist(sshCfgFile) || isLink {
		if err = os.RemoveAll(sshCfgFile); err != nil {
			return errors.Wrap(err, "remove ~/.ssh/config error")
		}
	}
	if err = os.Symlink(ctxSSHCfgFile, sshCfgFile); err != nil {
		return err
	}

	if err, isLinkdir = utils.IsSymLink(filepath.Join(env.SSHPath, "include")); err != nil {
		// return err
		isLinkdir = false
	}

	if utils.PathExist(filepath.Join(env.SSHPath, "include")) || isLinkdir {
		err = os.RemoveAll(filepath.Join(env.SSHPath, "include"))
		if err != nil {
			return errors.Wrap(err, "remove ~/.ssh include dir error")
		}
	}

	err = os.Symlink(filepath.Join(env.StorePath, "contexts", ctxName, "include"), filepath.Join(env.SSHPath, "include"))
	if err != nil {
		return err
	}
	return nil
}
