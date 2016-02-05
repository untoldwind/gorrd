package commands

import (
	"github.com/codegangsta/cli"
)

// Command line create command
var UpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update values to an rrd rrd file",
	Flags:     []cli.Flag{},
	ArgsUsage: "file.rrd",
	Action:    updateCommand,
}

func updateCommand(ctx *cli.Context) {
	_, err := getFilenameArg(ctx)
	if err != nil {
		showError(ctx, err)
		return
	}

}
