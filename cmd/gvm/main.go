package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gizmo-ds/gvm/cmd/gvm/commands"
	"github.com/urfave/cli/v2"
)

var AppVersion = "development"

func main() {
	_ = (&cli.App{
		Name:     "gvm",
		Usage:    "Golang Version Manager",
		Version:  AppVersion,
		Suggest:  true,
		Commands: commands.Commands,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "gvm-dir",
				Usage:   "The root directory of gvm installations",
				EnvVars: []string{color.GreenString("env:GVM_DIR")},
				Action: func(ctx *cli.Context, s string) error {
					return os.Setenv("GVM_DIR", s)
				},
			},
			&cli.BoolFlag{
				Name:               "no-color",
				Usage:              "Disable colored output",
				Value:              false,
				DisableDefaultText: true,
				EnvVars:            []string{color.GreenString("env:GVM_NO_COLOR")},
				Action: func(ctx *cli.Context, b bool) error {
					color.NoColor = b
					return os.Setenv("GVM_NO_COLOR", func() string {
						if b {
							return "true"
						}
						return "false"
					}())
				},
			},
			&cli.StringFlag{
				Name:    "mirror",
				Usage:   "Mirror of https://go.dev/dl/",
				Aliases: []string{"m"},
				EnvVars: []string{color.GreenString("env:GVM_DL_MIRROR")},
				Action: func(ctx *cli.Context, s string) error {
					return os.Setenv("GVM_DL_MIRROR", s)
				},
			},
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}).Run(os.Args)
}
