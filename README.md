# myssh
ssh 管理工具

主要是将[mmh](https://github.com/mritd/mmh)、[skm](https://github.com/TimothyYe/skm)、[manssh](https://github.com/xwjdsh/manssh)这三个ssh工具的功能整合到一起方便使用，部分代码拷贝自这三个工具。

## usage

```
MYSSH V0.0.1
https://github.com/cnlubo/myssh

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
   部分命令被重命名为快捷命令方便操作(安装后自动软连接)

- myssh alias ==> malias
- myssh cfg ==> mcfg
- myssh km  ==> mkm
- myssh clusters ==> mclusters


## SSH Keys Manager

```
mkm
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






