package myssh

import (
	"context"
	"errors"
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	cancelChannel := make(chan os.Signal, 1)
	signal.Notify(cancelChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGTSTP,
	)
	go func() {
		// read os exit signal
		s := <-cancelChannel
		fmt.Println()
		fmt.Println("退出信号:", s)
		// exit all goroutine
		cancel()

	}()
	if len(servers) == 1 {
		// single server
		err := execCommands(ctx, servers[0], true, cmd)
		utils.PrintErr(err)
	} else {

		{
			// multiple servers
			// create goroutine
			var execWg sync.WaitGroup
			execWg.Add(len(servers))
			for _, s := range servers {
				// async exec
				// because it takes time for ssh to establish a connection
				go func(s *ServerConfig) {
					defer execWg.Done()
					err := execCommands(ctx, s, true, cmd)
					utils.PrintErrWithPrefix(s.Name+": 😱 ", err)
				}(s)
			}
			execWg.Wait()
		}
	}
	return nil
}

// single server execution command
// since multiple tasks are executed async, the error is returned using channel
func execCommands(ctx context.Context, s *ServerConfig, singleServer bool, cmd string) error {
	// get ssh client
	sshClient, err := s.wrapperClient()
	if err != nil {
		// errCh <- err
		return err
	}
	defer func() {
		_ = sshClient.Close()
	}()

	// get ssh session
	session, err := s.wrapperSession(sshClient)
	if err != nil {
		// errCh <- err
		return err
	}

	// ssh utils session
	sshSession := utils.NewSSHSession(session)
	defer func() { _ = sshSession.Close() }()
	go func() {
		select {
		case <-ctx.Done():
			_ = sshSession.Close()
			_ = sshClient.Close()
		}
	}()

	return sshSession.PipeExec(cmd, func(r io.Reader, w io.Writer) {
		utils.Converted2Rendered(r, w, s.Name+":")
	})
}