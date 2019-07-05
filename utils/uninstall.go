package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
)

func Uninstall(osName string, dir string) error {

	if len(dir) == 0 {
		return errors.New("uninstall dir set failed!")
	}
	binPath := "/usr/local/bin"
	var binPaths = []string{
		filepath.Join(binPath, "mex"),    // exec
		filepath.Join(binPath, "mgo"),    // go
		filepath.Join(binPath, "msrv"),   // server
		filepath.Join(binPath, "mcfg"),   // cfg
		filepath.Join(binPath, "mkm"),    // km
		filepath.Join(binPath, "myssh"),  // myssh
		filepath.Join(binPath, "malias"), // host config
	}

	currentPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil
	}

	if osName == "linux" && !Root() {
		// cmd := exec.Command("sudo", currentPath, "uninstall", "--dir", dir)
		cmd := exec.Command("sudo", currentPath, "uninstall")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {

		for _, bin := range binPaths {
			if PathExist(bin) {
				PrintN(Uinst, fmt.Sprintf("remove %s\n", bin))
				_ = os.Remove(bin)
			}
		}
		if PathExist(filepath.Join(dir, "myssh")) {
			PrintN(Uinst, fmt.Sprintf("remove %s\n", filepath.Join(dir, "myssh")))
			_ = os.Remove(filepath.Join(dir, "myssh"))
			_ = os.RemoveAll(dir)
		}
	}
	return nil
}
