package myssh

import (
	"fmt"
	"github.com/cnlubo/myssh/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func BackupAll(backupPath string, env *Environment) error {

	backupDir := filepath.Join(backupPath, GetBakDirName("myssh"))
	if exists := utils.PathExist(backupDir); !exists {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			return errors.Wrap(err, "create backup dir failed")
		}
	}

	utils.PrintN(utils.Info, fmt.Sprintf("backup dir ==> [%s] \n", backupDir))

	// backup ./ssh/config file
	sshConfigFile := filepath.Join(env.SSHPath, "config")
	dstFile := filepath.Join(backupDir, "ssh_config")
	result := utils.Execute(env.SSHPath, "cp", sshConfigFile, dstFile)
	if result {
		utils.PrintN(utils.Info, "backup ssh config successfully\n")
	}
	// backup storePath
	destDir := filepath.Join(backupDir, "myssh")
	result = utils.Execute(env.StorePath, "cp", "-r", env.StorePath+"/.", destDir)
	if result {
		utils.PrintN(utils.Info, fmt.Sprintf("backup [%s] successfully\n", env.StorePath))
	}
	// backup skm
	destDir = filepath.Join(backupDir, "skm")
	result = utils.Execute(env.SKMPath, "cp", "-r", env.SKMPath+"/.", destDir)

	if result {
		utils.PrintN(utils.Info, fmt.Sprintf("backup [%s] successfully\n", env.SKMPath))
	}
	return nil
}
