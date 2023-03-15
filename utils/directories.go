package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func GvmDir() (string, error) {
	dir := os.Getenv("GVM_DIR")
	if dir != "" {
		return dir, nil
	}

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}
		dir = filepath.Join(dir, "gvm")

	case "darwin":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir = filepath.Join(dir, "Library", "Application Support", "gvm")

	default: // Unix
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_DATA_HOME nor $HOME are defined")
			}
			dir = filepath.Join(dir, ".gvm")
		} else {
			dir = filepath.Join(dir, "gvm")
		}
	}
	return dir, nil
}

func Mkdir(d string, e ...error) (string, error) {
	if len(e) > 0 && e[0] != nil {
		return d, e[0]
	}
	if _, err := os.Stat(d); os.IsNotExist(err) {
		return d, os.MkdirAll(d, os.ModePerm)
	}
	return d, nil
}
