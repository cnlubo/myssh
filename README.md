<!--
 * @Author: cnak47
 * @Date: 2019-08-20 22:17:29
 * @LastEditors: cnak47
 * @LastEditTime: 2019-08-20 22:45:34
 * @Description: 
 -->

# myssh

ssh 管理工具

主要是将[sshbatch](https://github.com/agentzh/sshbatch)、[mmh](https://github.com/mritd/mmh)、[skm](https://github.com/TimothyYe/skm)、[manssh](https://github.com/xwjdsh/manssh)这几个工具的功能整合到一起方便使用，部分代码拷贝自这几个工具。

## Usage

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
mkm -h
Manage multiple SSH key.

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

## manage myssh configfile

```bash
mcfg
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

## manage clusters

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

## managing your ssh alias config

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
