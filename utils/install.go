package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(osName string, dir string) error {

	if len(dir) == 0 {
		return errors.New("install dir set failed!")
	}

	installDir := filepath.Join(dir, "myssh")
	binPath := "/usr/local/bin"
	var binPaths = []string{
		filepath.Join(binPath, "mcfg"),      // cfg
		filepath.Join(binPath, "mkm"),       // km
		filepath.Join(binPath, "malias"),    // ssh alias
		filepath.Join(binPath, "mclusters"), // host cluster
		filepath.Join(binPath, "myssh"),     // myssh

	}
	// 获得当前文件执行路径和名称
	currentPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	if osName == "linux" && !Root() {
		// cmd := exec.Command("sudo", currentPath, "install", "--dir", dir)
		cmd := exec.Command("sudo", currentPath, "install")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "sudo run myssh failed")
		}
	} else {

		err := Uninstall(osName, installDir)
		// create install dir
		if PathExist(installDir) {
			if ok, _ := IsEmpty(installDir); ok {
				err := os.Remove(installDir)
				if err != nil {
					return errors.Wrap(err, "Remove installDir failed")
				}
			} else {
				return errors.New(fmt.Sprintf("install dir [%s] not empty", installDir))
			}
		}
		err = os.MkdirAll(installDir, 0755)
		if err != nil {
			return errors.Wrap(err, "Create installDir failed")
		}

		f, err := os.Open(currentPath)
		if err != nil {
			return errors.Wrap(err, "open Path failed")
		}

		defer func() {
			_ = f.Close()
		}()
		target, err := os.OpenFile(filepath.Join(installDir, "myssh"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return errors.Wrap(err, "OpenFile failed")
		}
		defer func() {
			_ = target.Close()
		}()
		PrintN(Inst, "install myssh.........\n")
		PrintN(Inst, fmt.Sprintf("install %s\n", filepath.Join(installDir, "myssh")))

		_, err = io.Copy(target, f)
		if err != nil {
			return errors.Wrap(err, "copy failed")
		}
		// create link
		for _, bin := range binPaths {
			PrintN(Inst, fmt.Sprintf("install link %s\n", bin))
			err = os.Symlink(filepath.Join(installDir, "myssh"), bin)
			if err != nil {
				return errors.Wrap(err, "create link failed")
			}

		}
	}
	return nil
}
