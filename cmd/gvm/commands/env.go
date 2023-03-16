package commands

import (
	"fmt"
	"strings"

	"github.com/gizmo-ds/gvm/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands, Env)
}

var env utils.GvmEnv
var Env = &cli.Command{
	Name:            "env",
	Usage:           "Print the GVM environment variables",
	HideHelpCommand: true,
	ArgsUsage:       " ",
	Before: func(ctx *cli.Context) error {
		env = utils.GvmEnv{}
		env.Init()
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "shell",
			Usage:    "The shell syntax to use. Print all when missing\noptional shell: bash, zsh, fish, powershell",
			Value:    "bash",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		if !ctx.IsSet("shell") {
			return fmt.Errorf("missing shell flag")
		}
		shell := ctx.String("shell")
		var exportStatement string
		switch shell {
		case "bash", "zsh":
			exportStatement = "export %s=\"%s\"\n"
		case "fish":
			exportStatement = "set -gx %s \"%s\";\n"
		case "powershell":
			exportStatement = "$env:%s=\"%s\"\n"
		default:
			return fmt.Errorf("unsupported shell: %s", shell)
		}
		var builder strings.Builder
		for _, env := range env.Envs() {
			builder.WriteString(fmt.Sprintf(exportStatement, env.Name, strings.ReplaceAll(env.Value, "\\", "\\\\")))
		}
		if shell == "zsh" {
			builder.WriteString("rehash\n")
		}
		fmt.Print(builder.String())
		return nil
	},
}
