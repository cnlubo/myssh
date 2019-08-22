<!--
 * @Author: cnak47
 * @Date: 2019-08-20 22:17:29
 * @LastEditors: cnak47
 * @LastEditTime: 2019-08-21 16:19:37
 * @Description: 
 -->

# myssh

使用 Go 编写的服务器SSH管理工具,主要将[sshbatch](https://github.com/agentzh/sshbatch)、[mmh](https://github.com/mritd/mmh)、[skm](https://github.com/TimothyYe/skm)、[manssh](https://github.com/xwjdsh/manssh)这几个工具的功能整合到一起以方便使用,部分代码拷贝自这几个工具。

## 安装

可直接从 [release](https://github.com/cnlubo/myssh/releases) 页下载预编译的二进制文件，然后执行 myssh install 安装,卸载直接执行 myssh uninstall,卸载命令不会删除 ~/.myssh 配置目录。

## 基本命令

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

所有的 SSH keys 都存储在 \$HOME/.mkm目录中,如果$HOME/.ssh目录下存在 id_rsa & id_rsa.pub key pairs 将被移动到\$HOME/.mkm/default.

### 创建 SSH key

当前支持 RSA 和 ED25519 两种类型的SSH key

## 管理配置文件

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

## 管理 SSH alias config

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
### 创建 SSH alias

```bash
malias add
```

![malias-add](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-add.gif)

### 查询 SSH alias

![malias-list](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-list.gif)

支持模糊查询，例如:malias ls test

### 删除 one or more alias

```bash
malias del test test-1
✔ deleted successfully!!!
```
### 修改 SSH alias

```bash
malias update test3
```
![malias-update](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-update.gif)


### 拷贝 SSH public key to a alias Host

```bash
Copy SSH public key to alias Host

Usage:
  myssh alias keycopy alias...

Aliases:
  keycopy, kcp

Flags:
  -h, --help   help for keycopy

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")
```

![malias-kcp](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-kcp.gif)

### 快捷登录服务器

```bash
malias go test1
```
### 批量执行命令

```bash
malias bt -h
Batch exec command for alias.

Usage:
  myssh alias batch alias command ... [flags]

Aliases:
  batch, bt

Flags:
  -h, --help     help for batch
  -P, --prompt   Prompt for password

Global Flags:
      --configPath string   Path where store myssh profiles.
                            can also be set by the MYSSH_CONFIG_HOME environment variable. (default "/Users/ak47/.myssh")
      --mkmPath string      Path where myssh should store multi SSHKeys.
                            can also be set by the MKM_PATH environment variable. (default "/Users/ak47/.mkm")
      --no-color            Disable color when outputting message.
      --sshPath string      Path to .ssh folder.
                            can also be set by the SSH_PATH environment variable. (default "/Users/ak47/.ssh")
```
![malias-batch](https://github.com/cnlubo/myssh/blob/master/snapshots/malias-batch.gif)


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
### 解析 Host pattern to machine host list

参考 [sshbatch](https://github.com/agentzh/sshbatch)中的fornodes脚本的使用方法，支持集合操作.

```bash
mclusters ls

   CLUSTERNAME              HOSTPATTERN
 ---------------- --------------------------------
  A                foo[01-03].com bar.org
  B                bar.org baz[a-b,d,e-g].cn
                   foo02.com
  C                {A} * {B}
  D                {A} - {B}

```
![mclusters](https://github.com/cnlubo/myssh/blob/master/snapshots/mclusters-1.gif)

