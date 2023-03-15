package gvm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gizmo-ds/gvm/utils"
)

func Install(version string, archive string) error {
	if info, err := os.Stat(archive); err != nil {
		return err
	} else if info.IsDir() {
		return fmt.Errorf("%s is a directory", archive)
	}
	dir, err := utils.GvmDir()
	if err != nil {
		return err
	}
	dest, err := utils.Mkdir(filepath.Join(dir, "versions", version))
	if err != nil {
		return err
	}
	return utils.Extract(archive, dest)
}

func Uninstall(version string) error {
	dir, err := utils.GvmDir()
	if err != nil {
		return err
	}
	dir = filepath.Join(dir, "versions", version)
	return os.RemoveAll(dir)
}
