<!--
 * @Author: cnak47
 * @Date: 2019-08-20 22:17:29
 * @LastEditors: cnak47
 * @LastEditTime: 2019-08-21 16:19:37
 * @Description: 
 -->

# myssh

ssh 管理工具

将[sshbatch](https://github.com/agentzh/sshbatch)、[mmh](https://github.com/mritd/mmh)、[skm](https://github.com/TimothyYe/skm)、[manssh](https://github.com/xwjdsh/manssh)这几个工具的功能整合到一起以方便使用，部分代码拷贝自这几个工具。

## 使用

```bash
myssh -h

My ssh toolkit. Flags and arguments can be input to do what actually you wish.

Usage:
  myssh [flags]
  myssh [command]

Available Commands:
  km          manage ssh key (alias: mkm)
  cfg         manage myssh configfile (alias: mcfg)
  clusters    manage clusters (alias: mclusters)
  alias       managing your ssh alias config (alias: malias)
  install     install myssh
  uninstall   uninstall myssh
  backup      Backup all config,SSHKeys...
  version     Print versions about Myssh

Flags:
      --no-color            Disable color when outputting message.
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")
  -h, --help                help for myssh

Use "myssh [command] --help" for more information about a command.

```

部分命令安装后自动软连接为快捷命令以便方便使用。

- myssh alias ==> malias
- myssh cfg ==> mcfg
- myssh km  ==> mkm
- myssh clusters ==> mclusters

## SSH Keys Manager

管理 多个 SSH Keys

```bash
Manage multiple SSH keys

Usage:
  myssh km
  myssh km [command]

Aliases:
  km, mkm

Available Commands:
  init        Initialize SSH keys store
  list        List all available SSH keys (alias:ls)
  add         add one SSHKey
  delete      delete SSH key (alias:del)
  use         Set specific SSH key as default
  display     Display SSHKey (alias:dp)
  rename      Rename SSH key aliasName (alias:rn)
  copy        Copy SSH public key to a remote host (alias:cp)

Flags:
  -h, --help   help for km

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")

Use "myssh km [command] --help" for more information about a command.

```

![mkm](https://github.com/cnlubo/myssh/blob/master/snapshots/mkm.gif)

### 初始化 SSH Key Store

第一次运行时需要初始化 SSH Key Store

```bash
mkm init
```

所有的 SSH key 都存储在 \$HOME/.mkm目录中,如果$HOME/.ssh目录下存在 id_rsa & id_rsa.pub key pairs 将被移动到\$HOME/.mkm/default.

### 创建 SSH key

Currently ONLY RSA and ED25519 keys are supported!

## Manage configfile

```bash
mcfg -h
manage myssh ConfigFile.

Usage:
  myssh cfg
  myssh cfg [command]

Aliases:
  cfg, mcfg

Available Commands:
  list        List config (alias:ls)
  add         Add context config
  delete      delete config
  set         set current config

Flags:
  -h, --help   help for cfg

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")

Use "myssh cfg [command] --help" for more information about a command.
```

![mcfg](https://github.com/cnlubo/myssh/blob/master/snapshots/mcfg.gif)

### 配置文件

支持多配置文件切换功能,安装完成后将会自动在$HOME/.myssh下创建默认配置，默认配置文件结构如下:

```bash
tree ~/.myssh
/Users/ak47/.myssh
├── contexts
│   ├── ak47
│   │   ├── ak47.yaml
│   │   ├── include
│   │   └── sshconfig
│   └── default
│       ├── default.yaml
│       ├── include
│       │   └── k8s
│       └── sshconfig
└── main.yaml
```

#### main.yaml

主配置文件结构如下:

```bash
basic: default
contexts:
- name: ak47
  sshconfig: /Users/ak47/.myssh/contexts/ak47/sshconfig
  clusterCfg: /Users/ak47/.myssh/contexts/ak47/ak47.yaml
- name: default
  sshconfig: /Users/ak47/.myssh/contexts/default/sshconfig
  clusterCfg: /Users/ak47/.myssh/contexts/default/default.yaml
current: default
```

配置文件中可以配置多个context,由current字段指明当前使用哪个context.

#### sshconfig

默认内容如下:

```bash
Include include/*
```

当前使用的context的 sshconfig 会被软连接到 ~/.ssh/config文件。可以使用 malias 命令 管理此配置文件。

#### clusterCfg

默认内容如下:

```bash
default:
  user: ak47
  privateKey: /Users/ak47/.ssh/id_ed25519
  port: 22
  server_alive_interval: 30s
clusters: []
```
文件名一般与context名称相同,可以使用mclusters 命令来管理此配置文件.

## 管理 ssh alias config 文件

```bash
malias -h
command line tool for managing your ssh alias config.

Usage:
  myssh alias
  myssh alias [command]

Aliases:
  alias, malias

Available Commands:
  list        List ssh alias (alias:ls)
  delete      Delete ssh alias (alias:del)
  add         Add one ssh alias
  update      Update ssh alias
  batch       batch exec command (alias: bt)
  go          ssh login Server
  keycopy     Copy SSH public key to a alias Host (alias:kcp)

Flags:
  -h, --help   help for alias

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")

Use "myssh alias [command] --help" for more information about a command.
```
### Add a new alias

```bash
malias add
```

![malias-add](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-add.gif)

### List or query alias

![malias-list](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-list.gif)




## manage clusters

整合了[sshbatch](https://github.com/agentzh/sshbatch)的相关功能,具体使用方法参考[sshbatch](https://github.com/agentzh/sshbatch).

```bash
mclusters -h
manage clusters.

Usage:
  myssh clusters
  myssh clusters [command]

Aliases:
  clusters, mclusters

Available Commands:
  parse       parse hostPattern to host list (alias:pa)
  add         Add one cluster
  list        List all cluster (alias:ls)
  delete      delete cluster (alias:del)
  batch       batch exec command (alias: bt)
  keycopy     Copy public key to cluster (alias: kcp)

Flags:
  -h, --help   help for clusters

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")

Use "myssh clusters [command] --help" for more information about a command.
```
### parse host pattern to host list

```bash
mclusters pa -h
parse host pattern to host list.

Usage:
  myssh clusters parse HOST_PATTERN [flag]

Aliases:
  parse, pa

Flags:
  -e, --expand   Expand the host list output to multiple lines
  -h, --help     help for parse

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")
```