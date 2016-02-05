package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/go-errors/errors"
)

func getFilenameArg(ctx *cli.Context) (string, error) {
	if !ctx.Args().Present() {
		return "", errors.Errorf("Filename required")
	}
	return ctx.Args().First(), nil
}

func showError(ctx *cli.Context, err error) {
	fmt.Fprintln(ctx.App.Writer, err.Error())
	fmt.Fprintln(ctx.App.Writer)

	cli.ShowCommandHelp(ctx, ctx.Command.Name)
}
