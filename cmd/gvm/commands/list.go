package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gizmo-ds/gvm"
	"github.com/gizmo-ds/gvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands, ListLocal, ListRemote)
}

var ListRemote = &cli.Command{
	Name:    "list-remote",
	Usage:   "List remote versions",
	Aliases: []string{"lr"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "Show all remote versions",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "stable",
			Aliases: []string{"s"},
			Usage:   "Show only stable versions",
			Value:   true,
		},
	},
	Action: func(ctx *cli.Context) error {
		stable := []bool{ctx.Bool("stable")}
		if ctx.Bool("all") {
			stable = nil
		}
		versions, err := gvm.ListRemote(stable...)
		if err != nil {
			return err
		}
		utils.Reverse(versions)

		for _, v := range versions {
			str := strings.TrimPrefix(v.Version, "go")
			if !v.Stable {
				str = color.YellowString("%s (unstable)", str)
			} else {
				str = color.GreenString(str)
			}
			fmt.Println(str)
		}
		return nil
	},
}

var ListLocal = &cli.Command{
	Name:    "list",
	Usage:   "List all locally installed golang versions",
	Aliases: []string{"ls"},
	Action: func(ctx *cli.Context) error {
		versions, _ := gvm.ListLocalVersions()
		currentVersion, isGvm, _ := gvm.CurrentVersion()
		currentVersion = strings.TrimPrefix(currentVersion, "go")

		system := "* system"
		if !isGvm {
			system += " -> " + currentVersion
			system = current(system)
		}
		fmt.Println(system)

		for _, v := range versions {
			version := strings.TrimPrefix(v, "go")
			str := fmt.Sprintf("* %s", version)
			if isGvm && version == currentVersion {
				str = current(str)
			}
			fmt.Println(str)
		}
		return nil
	},
}

func current(s string) string {
	return color.GreenString(s) + " (current)"
}
