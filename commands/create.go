package commands

import (
	"github.com/codegangsta/cli"
)

// Command line create command
var CreateCommand = cli.Command{
	Name:   "create",
	Usage:  "Create new rrd file",
	Flags:  []cli.Flag{},
	Action: createComamnd,
}

func createComamnd(ctx *cli.Context) {

}
