package gvm

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gizmo-ds/gvm/utils"
)

func CurrentVersion() (string, bool, error) {
	dir, err := utils.GvmDir()
	if err != nil {
		return "", false, err
	}
	dir = filepath.Join(dir, "versions")
	file := "go"
	if runtime.GOOS == "windows" {
		file += ".exe"
	}
	path, err := exec.LookPath(file)
	if err != nil {
		return "", false, err
	}
	version, err := goVersion(path)
	if err != nil {
		return "", false, err
	}
	versionDir := filepath.Join(filepath.Dir(path), "..")
	return version, strings.HasPrefix(versionDir, dir), err
}

func goVersion(file string) (string, error) {
	cmd := exec.Command(file, "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	str := strings.TrimSpace(out.String())
	arr := strings.Split(str, " ")
	for _, s := range arr {
		if len(s) > 2 && s[0:2] == "go" {
			return s, nil
		}
	}
	return "", errors.New("could not find go version")
}
