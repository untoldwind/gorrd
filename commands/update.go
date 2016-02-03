package commands

import (
	"github.com/codegangsta/cli"
)

// Command line create command
var UpdateCommand = cli.Command{
	Name:   "update",
	Usage:  "Update values to an rrd rrd file",
	Flags:  []cli.Flag{},
	Action: updateCommand,
}

func updateCommand(ctx *cli.Context) {

}
