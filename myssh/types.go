package myssh

import (
	"time"
)

type Environment struct {
	StorePath string
	SSHPath   string
	SKMPath   string
}

// ListOption options for List
type ListOption struct {
	// Keywords set Keyword filter records
	Keywords []string
	// IgnoreCase ignore case
	IgnoreCase bool
}

// server config
type ServerConfig struct {
	Name                string        `yaml:"name"`
	User                string        `yaml:"user"`
	Password            string        `yaml:"password"`
	Method              string        `yaml:"method"` // auth method default：key，options:password、key
	SuRoot              bool          `yaml:"su_root"`
	UseSudo             bool          `yaml:"use_sudo"`
	NoPasswordSudo      bool          `yaml:"no_password_sudo"`
	RootPassword        string        `yaml:"root_password"`
	PrivateKey          string        `yaml:"privateKey"`
	PrivateKeyPassword  string        `yaml:"privateKey_password"`
	Address             string        `yaml:"address"`
	Port                int           `yaml:"port"`
	ServerAliveInterval time.Duration `yaml:"server_alive_interval"`
}

// myssh servers
type Servers []*ServerConfig

func (servers Servers) Len() int {
	return len(servers)
}
func (servers Servers) Less(i, j int) bool {
	return servers[i].Name < servers[j].Name
}
func (servers Servers) Swap(i, j int) {
	servers[i], servers[j] = servers[j], servers[i]
}

// ************************************************************************

type Clusters []*ClusterConfig

func (cs Clusters) Len() int {
	return len(cs)
}
func (cs Clusters) Less(i, j int) bool {
	return cs[i].Name < cs[j].Name
}
func (cs Clusters) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

type ClusterConfig struct {
	Name        string `yaml:"name"`
	HostPattern string `yaml:"hostPattern"`
}

type ClustersConfig struct {
	configPath string
	Default    DefaultClusterConfig `yaml:"default"`
	Clusters   Clusters             `yaml:"clusters"`
}

type DefaultClusterConfig struct {
	User                string        `yaml:"user"`
	PrivateKey          string        `yaml:"privateKey"`
	Port                int           `yaml:"port"`
	ServerAliveInterval time.Duration `yaml:"server_alive_interval"`
}

func (cfg *ClustersConfig) SetConfigPath(configPath string) {
	cfg.configPath = configPath
}

func (cfg *ClustersConfig) ConfigPath() string {
	return cfg.configPath
}

// myssh context
type Context struct {
	Name          string `yaml:"name"`
	SSHConfig     string `yaml:"sshconfig"`
	ClusterConfig string `yaml:"clusterCfg"`
}

// myssh contexts
type Contexts []Context

func (cs Contexts) Len() int {
	return len(cs)
}
func (cs Contexts) Less(i, j int) bool {
	return cs[i].Name < cs[j].Name
}
func (cs Contexts) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

// main config struct
type MainConfig struct {
	configPath string
	Basic      string   `yaml:"basic"`
	Contexts   Contexts `yaml:"contexts"`
	Current    string   `yaml:"current"`
}
