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
				Name:  "gvm-dir",
				Usage: "The root directory of gvm installations",
				Action: func(ctx *cli.Context, s string) error {
					return os.Setenv("GVM_DIR", s)
				},
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable colored output. Use --no-color to disable.",
				Value: false,
				Action: func(ctx *cli.Context, b bool) error {
					color.NoColor = b
					return nil
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
