package myssh

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"github.com/gookit/color"
	"github.com/mritd/sshutils"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"
)

func batchExec(servers Servers, cmd string) error {

	if len(servers) == 0 {
		return errors.New("servers not exists")
	}
	// clear screen
	//  utils.Clear()
	// use context to manage goroutine
	ctx, cancel := context.WithCancel(context.Background())

	// monitor os signal
	cancelChannel := make(chan os.Signal)
	signal.Notify(cancelChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		switch <-cancelChannel {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			// exit all goroutine
			cancel()
		}
	}()
	if len(servers) == 1 {
		// single server
		var errCh = make(chan error, 1)
		execCommands(ctx, servers[0], true, cmd, errCh)
		select {
		case err := <-errCh:
			utils.PrintErr(err)
			fmt.Println()
		default:
		}
	} else {
		// create goroutine
		var serverWg sync.WaitGroup
		serverWg.Add(len(servers))
		for _, tmpServer := range servers {
			server := tmpServer
			// async exec
			// because it takes time for ssh to establish a connection
			go func() {
				defer serverWg.Done()
				var errCh = make(chan error, 1)
				execCommands(ctx, server, false, cmd, errCh)
				select {
				case err := <-errCh:
					utils.PrintErr(errors.Wrap(err, server.Name))
					fmt.Println()
				default:
				}
			}()
		}
		serverWg.Wait()
	}
	return nil
}

// single server execution command
// since multiple tasks are executed async, the error is returned using channel
func execCommands(ctx context.Context, s *ServerConfig, singleServer bool, cmd string, errCh chan error) {

	// get ssh client
	sshClient, err := s.sshClient()
	if err != nil {
		errCh <- err
		return
	}
	defer func() {
		_ = sshClient.Close()
	}()

	// get ssh session
	session, err := sshClient.NewSession()
	if err != nil {
		errCh <- err
		return
	}

	// ssh utils session
	sshSession := sshutils.NewSSHSession(session)
	defer func() {
		_ = sshSession.Close()
	}()

	// exec cmd
	go sshSession.PipeExec(cmd)

	// copy error
	var errWg sync.WaitGroup
	errWg.Add(1)
	go func() {
		// ensure that the error message is successfully output
		defer errWg.Done()
		select {
		case err, ok := <-sshSession.Error():
			if ok {
				errCh <- err
			}
		}
	}()

	// print to stdout
	go func() {
		select {
		case <-sshSession.Ready():
			// read from sshSession.Stdout and print to os.stdout
			if singleServer {
				_, _ = io.Copy(os.Stdout, sshSession.Stdout)
			} else {
				f, _ := utils.GetColor()
				t, err := template.New("").Funcs(utils.ColorFuncMap).Parse(fmt.Sprintf(`{{ .Name | %s}}{{ ":" | %s}}  {{ .Value}}`, f, f))

				if err != nil {
					errCh <- err
					return
				}

				buf := bufio.NewReader(sshSession.Stdout)
				for {
					line, err := buf.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							break
						} else {
							errCh <- err
							break
						}
					}
					var output bytes.Buffer
					err = t.Execute(&output, struct {
						Name  string
						Value string
					}{
						Name:  s.Name,
						Value: color.FgLightWhite.Render(string(line)),
					})
					if err != nil {
						errCh <- err
						break
					}
					fmt.Print(output.String())
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		_ = sshClient.Close()
	case <-sshSession.Done():
		_ = sshClient.Close()
	}

	errWg.Wait()

}
