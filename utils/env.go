package utils

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type GvmEnv struct {
	Dir      string `env:"GVM_DIR"`
	DlMirror string `env:"GVM_DL_MIRROR"`
	NoColor  bool   `env:"GVM_NO_COLOR" `
}

func (e *GvmEnv) Init() error {
	t := reflect.TypeOf(e).Elem()
	v := reflect.ValueOf(e).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		env := f.Tag.Get("env")
		switch f.Name {
		case "Dir":
			dir := os.Getenv(env)
			if dir == "" {
				var err error
				dir, err = GvmDir()
				if err != nil {
					return err
				}
			}
			_ = os.Setenv(env, dir)
			v.FieldByName(f.Name).SetString(dir)
		case "DlMirror":
			mirror := os.Getenv(env)
			switch strings.ToLower(mirror) {
			case "":
				mirror = "https://go.dev/dl/"
			case "china":
				mirror = "https://golang.google.cn/dl/"
			}
			_ = os.Setenv(env, mirror)
			v.FieldByName(f.Name).SetString(mirror)
		case "NoColor":
			noColor, ok := os.LookupEnv(env)
			if ok {
				v.FieldByName(f.Name).SetBool(noColor == "true")
			}
			_ = os.Setenv(env, noColor)
		}
	}
	return nil
}

type Env struct {
	Name  string
	Value string
}

func (e *GvmEnv) Envs() []Env {
	var envs []Env
	t := reflect.TypeOf(e).Elem()
	v := reflect.ValueOf(e).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		envs = append(envs, Env{
			Name:  f.Tag.Get("env"),
			Value: fmt.Sprint(v.FieldByName(f.Name)),
		})
	}
	return envs
}
