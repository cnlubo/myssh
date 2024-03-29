package myssh

import (
	"fmt"
	"github.com/cnlubo/myssh/utils"
	//"github.com/mritd/sshutils"

	//"github.com/mritd/sshutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strings"
	"time"
)

// set server default config
func (s *ServerConfig) setDefault() {

	if s.User == "" {
		s.User = utils.GetUsername()
	}
	if s.Port == 0 {
		s.Port = 22
	}

	if s.ServerAliveInterval == 0 {
		s.ServerAliveInterval = 30 * time.Second
	}
}

// func (s *ServerConfig) setDefault() {
// 	if s.User == "" {
// 		s.User = ContextCfg.Basic.User
// 	}
// 	if s.Port == 0 {
// 		s.Port = ContextCfg.Basic.Port
// 	}
// 	if s.Password == "" {
// 		// s.Password = ContextCfg.Basic.Password
// 		if s.PrivateKey == "" {
// 			s.PrivateKey = ContextCfg.Basic.PrivateKey
// 		} else {
// 			if ok := path.IsAbs(s.PrivateKey); !ok {
// 				home, _ := homedir.Dir()
// 				s.PrivateKey = filepath.Join(home+"/.ssh", filepath.Base(s.PrivateKey))
// 			}
// 		}
// 	}
// 	// if s.PrivateKeyPassword == "" {
// 	// 	s.PrivateKeyPassword = ContextCfg.Basic.PrivateKeyPassword
// 	// }
// 	// if s.Proxy == "" {
// 	// 	s.Proxy = ContextCfg.Basic.Proxy
// 	// }
// 	if s.ServerAliveInterval == 0 {
// 		s.ServerAliveInterval = ContextCfg.Basic.ServerAliveInterval
// 	}
// }

// wrapperClient return a standard ssh client with specific parameters set
func (s *ServerConfig) wrapperClient() (*ssh.Client, error) {

	return s.sshClient()
}

// return a ssh client intense point
func (s *ServerConfig) sshClient() (*ssh.Client, error) {

	// default to basic config
	s.setDefault()

	var client *ssh.Client
	auths, err := s.parseAuthMethods()
	if err != nil {
		return nil, errors.Wrap(err, "auth fail")
	}
	sshConfig := &ssh.ClientConfig{
		User:            s.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	client, err = ssh.Dial("tcp", fmt.Sprint(s.Address, ":", s.Port), sshConfig)

	if err != nil {
		if utils.ErrorAssert(err, "ssh: unable to authenticate") {
			return nil, errors.New("connect errors,please check privateKey or password")
		} else {
			return nil, errors.Wrap(err, "connect errors")
		}
	}

	return client, nil
}

// wrapperSession returns a standard ssh session with specific parameters set
// if there is an error, the ssh session is nil
func (s *ServerConfig) wrapperSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	//_ = session.Setenv(envMMH, "true")
	//if s.Environment != nil {
	//	for k, v := range s.Environment {
	//		_ = session.Setenv(k, v)
	//	}
	//}
	return session, nil
}

// start a ssh terminal
func (s *ServerConfig) Terminal() error {
	sshClient, err := s.sshClient()
	if err != nil {
		return err
	}
	defer func() {
		_ = sshClient.Close()
	}()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}

	var sshSession *utils.SSHSession
	if s.SuRoot {
		sshSession = utils.NewSSHSessionWithRoot(session, s.UseSudo, s.NoPasswordSudo, s.RootPassword, s.Password)
	} else {
		sshSession = utils.NewSSHSession(session)
	}

	defer func() {
		_ = sshSession.Close()
	}()

	// keep alive
	if s.ServerAliveInterval > 0 {
		return sshSession.TerminalWithKeepAlive(s.ServerAliveInterval)
	}
	return sshSession.Terminal()

}

// get auth method
// priority use privateKey method
func (s *ServerConfig) parseAuthMethods() ([]ssh.AuthMethod, error) {
	var sshs []ssh.AuthMethod
	if s.Method == "" || s.Method == "key" {
		if strings.TrimSpace(s.PrivateKey) != "" {
			method, err := privateKeyFile(s.PrivateKey, s.PrivateKeyPassword)
			if err != nil {
				return nil, err
			}
			sshs = append(sshs, method)
		} else {
			return nil, errors.New("empty privateKey")
		}
	} else {
		sshs = append(sshs, ssh.Password(s.Password))
	}

	return sshs, nil
}

// use private key to return ssh auth method
func privateKeyFile(file, password string) (ssh.AuthMethod, error) {

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer

	if strings.TrimSpace(password) == "" {
		signer, err = ssh.ParsePrivateKey(buffer)
	} else {
		fmt.Printf(password)
		signer, err = ssh.ParsePrivateKeyWithPassphrase(buffer, []byte(password))
	}

	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}