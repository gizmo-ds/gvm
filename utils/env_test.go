package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGvmEnv_Init(t *testing.T) {
	t.Run("GVM_DIR", func(t *testing.T) {
		e := &GvmEnv{}
		v := "/home/gizmo/.gvm"
		_ = os.Setenv("GVM_DIR", v)
		defer os.Unsetenv("GVM_DIR")
		err := e.Init()
		assert.NoError(t, err)
		assert.Equal(t, v, e.Dir)
	})
	t.Run("GVM_DL_MIRROR", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			e := &GvmEnv{}
			_ = os.Setenv("GVM_DL_MIRROR", "")
			defer os.Unsetenv("GVM_DL_MIRROR")
			err := e.Init()
			assert.NoError(t, err)
			assert.Equal(t, "https://go.dev/dl/", e.DlMirror)
		})
		t.Run("China", func(t *testing.T) {
			e := &GvmEnv{}
			_ = os.Setenv("GVM_DL_MIRROR", "china")
			defer os.Unsetenv("GVM_DL_MIRROR")
			err := e.Init()
			assert.NoError(t, err)
			assert.Equal(t, "https://golang.google.cn/dl/", e.DlMirror)
		})
		t.Run("Custom", func(t *testing.T) {
			e := &GvmEnv{}
			_ = os.Setenv("GVM_DL_MIRROR", "https://google.com")
			defer os.Unsetenv("GVM_DL_MIRROR")
			err := e.Init()
			assert.NoError(t, err)
			assert.Equal(t, "https://google.com", e.DlMirror)
		})
	})
	t.Run("GVM_NO_COLOR", func(t *testing.T) {
		e := &GvmEnv{}
		_ = os.Setenv("GVM_NO_COLOR", "true")
		defer os.Unsetenv("GVM_NO_COLOR")
		err := e.Init()
		assert.NoError(t, err)
		assert.Equal(t, true, e.NoColor)
	})
}

func TestGvmEnv_Envs(t *testing.T) {
	_ = os.Setenv("GVM_NO_COLOR", "true")
	defer os.Unsetenv("GVM_NO_COLOR")
	_ = os.Setenv("GVM_DL_MIRROR", "china")
	defer os.Unsetenv("GVM_DL_MIRROR")
	_ = os.Setenv("GVM_DIR", "/home/gizmo/.gvm")
	defer os.Unsetenv("GVM_DIR")

	e := &GvmEnv{}
	err := e.Init()
	assert.NoError(t, err)
	envs := e.Envs()
	for _, env := range envs {
		switch env.Name {
		case "GVM_NO_COLOR":
			assert.Equal(t, true, e.NoColor)
		case "GVM_DL_MIRROR":
			assert.Equal(t, "https://golang.google.cn/dl/", e.DlMirror)
		case "GVM_DIR":
			assert.Equal(t, "/home/gizmo/.gvm", e.Dir)
		}
	}
}
