package commands

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/gizmo-ds/gvm"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands, Install, Uninstall)
}

var Install = &cli.Command{
	Name:    "install",
	Usage:   "Install a version of Go",
	Aliases: []string{"i", "add"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "latest",
			Aliases: []string{"l"},
			Usage:   "Install the latest stable version",
		},
	},
	Action: func(ctx *cli.Context) error {
		versionStr := ctx.Args().First()
		if ctx.Bool("latest") {
			version, err := gvm.LatestStableVersion()
			if err != nil {
				return err
			}
			versionStr = version.Version
		}
		if versionStr == "" {
			return errors.New("no version specified")
		}
		versionStr = strings.TrimPrefix(versionStr, "go")

		fmt.Printf("Installing %s (%s)\n",
			color.BlueString(versionStr),
			runtime.GOARCH,
		)

		bar := pb.Default.New(0).
			Set(pb.Bytes, true).
			Set(pb.SIBytesPrefix, true)
		file, err := gvm.Download("go"+versionStr, os.TempDir(), bar)
		if err != nil {
			return err
		}
		if err = gvm.Install(versionStr, file); err != nil {
			return err
		}
		fmt.Printf("Successfully installed %s\n", color.GreenString(versionStr))
		return nil
	},
}

var Uninstall = &cli.Command{
	Name:    "uninstall",
	Usage:   "Uninstall a version of Go",
	Aliases: []string{"rm", "remove"},
	Action: func(ctx *cli.Context) error {
		versionStr := ctx.Args().First()
		if versionStr == "" {
			return errors.New("no version specified")
		}
		versionStr = strings.TrimPrefix(versionStr, "go")
		fmt.Printf("Uninstalling %s (%s)\n",
			color.BlueString(versionStr),
			runtime.GOARCH,
		)
		if err := gvm.Uninstall(versionStr); err != nil {
			return err
		}
		fmt.Printf("Successfully uninstalled %s\n", color.GreenString(versionStr))
		return nil
	},
}
