package myssh

import (
	"github.com/cnlubo/myssh/utils"
	"github.com/cnlubo/promptx"
	_ "github.com/gookit/color"
	//_ "github.com/mritd/sshutils"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

func ClusterCopy(args []string, PromptPass bool, login *ServerConfig) error {

	hostList, err := ParseExpr(args[0])
	if err != nil {
		return errors.Wrap(err, "bad hostPattern")
	}
	var servers Servers
	reg := regexp.MustCompile(`^(\w+)[@](.+$)`) // user@host:port
	for _, value := range hostList {
		var sshUser, loginPass, host, port, privateKey, authMethod string

		if reg.MatchString(utils.CompressStr(value)) {
			sshUser, host, port = utils.ParseConnect(utils.CompressStr(value))

		} else {
			sshUser = login.User
			host = utils.CompressStr(value)
			port = strconv.Itoa(login.Port)
		}
		pp, _ := strconv.Atoi(port)

		// set default
		if sshUser == "" {
			sshUser = ClustersCfg.Default.User
		}
		if pp == 0 {
			pp = ClustersCfg.Default.Port
		}
		if login.PrivateKey == "" {
			privateKey = ClustersCfg.Default.PrivateKey
		} else {
			privateKey = login.PrivateKey
		}
		// prompt pass
		if PromptPass {
			utils.PrintN(utils.Warn, "Server "+sshUser+"@"+host)
			loginPass = ""

			prompt := promptx.NewDefaultPrompt(func(line []rune) error {
				if strings.TrimSpace(string(line)) == "" {
					return inputEmptyErr
				}
				return nil

			}, "Please type login password:")

			prompt.Mask = MaskPrompt
			loginPass, err = prompt.Run()
			if err != nil {
				return err
			}
			authMethod = "password"

		} else {
			loginPass = ""
			authMethod = "key"
		}
		server := &ServerConfig{
			Name:       host,
			User:       sshUser,
			Address:    host,
			Port:       pp,
			Method:     authMethod,
			PrivateKey: privateKey,
			// PrivateKeyPassword: login.PrivateKeyPassword,
			Password:            loginPass,
			ServerAliveInterval: ClustersCfg.Default.ServerAliveInterval,
		}
		servers = append(servers, server)

	}
	// exCmd := strings.Join(cmd, "&&")
	// err = batchExec(servers, exCmd)
	// if err != nil {
	// 	return err
	// }
	return nil

}
func runUpload(servers Servers, sourcePath string, targetPath string) error {
	if len(servers) == 0 {
		return errors.New("servers not exists")
	}

	return nil
}

// func runCopy(servers Servers,args []string, singleServer bool) error {
//
// 	if len(args) < 2 {
// 		return errors.New("parameter invalid")
// 	}
// 	if len(servers) == 0 {
// 		return errors.New("servers not exists")
// 	}
// 	// var servers Servers
// 	// download, eg: mcp test:~/file localPath
// 	// only single file/directory download is supported
// 	if len(strings.Split(args[0], ":")) == 2 && len(args) == 2 {
//
// 		// only single server is supported
// 		serverName := strings.Split(args[0], ":")[0]
// 		remotePath := strings.Split(args[0], ":")[1]
// 		localPath := args[1]
// 		s := FindServerByName(serverName)
// 		if s == nil {
// 			return errors.New("server not found")
// 		} else {
// 			client, err := s.sshClient()
// 			if err != nil {
// 				return err
// 			}
// 			defer func() {
// 				_ = client.Close()
// 			}()
// 			scpClient, err := sshutils.NewSCPClient(client)
// 			if err != nil {
// 				return err
// 			}
// 			return scpClient.CopyRemote2Local(remotePath, localPath)
// 		}
//
// 		// upload, eg: mcp localFile1 localFile2 localDir test:~
// 	} else if len(strings.Split(args[len(args)-1], ":")) == 2 {
//
// 		serverOrTag := strings.Split(args[len(args)-1], ":")[0]
// 		remotePath := strings.Split(args[len(args)-1], ":")[1]
//
// 		// single server copy
// 		if singleServer {
// 			s := findServerByName(serverOrTag)
// 			if s == nil {
// 				return errors.New("server not found")
// 			} else {
// 				client, err := s.sshClient()
// 				if err != nil {
// 					return err
// 				}
// 				defer func() {
// 					_ = client.Close()
// 				}()
// 				scpClient, err := sshutils.NewSCPClient(client)
// 				if err != nil {
// 					return err
// 				}
// 				allArg := args[:len(args)-1]
// 				allArg = append(allArg, remotePath)
// 				return scpClient.CopyLocal2Remote(allArg...)
// 			}
// 		} else {
// 			// multi server copy
// 			servers := findServersByTag(serverOrTag)
// 			if len(servers) == 0 {
// 				return errors.New("tagged server not found")
// 			}
//
// 			var wg sync.WaitGroup
// 			wg.Add(len(servers))
//
// 			for _, s := range servers {
// 				tmpServer := s
// 				go func() {
// 					defer wg.Done()
// 					client, err := tmpServer.sshClient()
// 					if err != nil {
// 						_, _ = color.New(color.BgRed, color.FgHiWhite).Printf("%s:  %s", tmpServer.Name, err)
// 						return
// 					}
// 					defer func() {
// 						_ = client.Close()
// 					}()
// 					scpClient, err := sshutils.NewSCPClient(client)
// 					if err != nil {
// 						_, _ = color.New(color.BgRed, color.FgHiWhite).Printf("%s:  %s", tmpServer.Name, err)
// 						return
// 					}
//
// 					allArg := args[:len(args)-1]
// 					allArg = append(allArg, remotePath)
// 					err = scpClient.CopyLocal2Remote(allArg...)
// 					if err != nil {
// 						_, _ = color.New(color.BgRed, color.FgHiWhite).Printf("%s:  %s", tmpServer.Name, err)
// 						return
// 					}
// 				}()
// 			}
//
// 			wg.Wait()
// 		}
//
// 	} else {
// 		return errors.New("unsupported mode")
// 	}
//
// 	return nil
// }