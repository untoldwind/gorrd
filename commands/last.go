package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/untoldwind/gorrd/rrd/cdata"
)

// Command line create command
var LastCommand = cli.Command{
	Name:      "last",
	Usage:     "Show the timestamp of the last update to a rrd file",
	Flags:     []cli.Flag{},
	ArgsUsage: "file.rrd",
	Action:    last,
}

func last(ctx *cli.Context) {
	args := ctx.Args()

	if !args.Present() {
		fmt.Fprintln(ctx.App.Writer, "Filename required")
		fmt.Fprintln(ctx.App.Writer)

		cli.ShowCommandHelp(ctx, "lastupdate")
		return
	}

	rrd, err := cdata.OpenRrdRawFile(args.First(), true)
	if err != nil {
		fmt.Fprintln(ctx.App.Writer, err)
		return
	}
	defer rrd.Close()

	fmt.Fprintf(ctx.App.Writer, "%d\n", rrd.LastUpdate.Unix())
}
