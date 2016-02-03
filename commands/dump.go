package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/untoldwind/gorrd/rrd/cdata"
)

// Command line create command
var DumpCommand = cli.Command{
	Name:      "dump",
	Usage:     "Dump contents of an rrd file",
	Flags:     []cli.Flag{},
	ArgsUsage: "file.rrd",
	Action:    dumpCommand,
}

func dumpCommand(ctx *cli.Context) {
	args := ctx.Args()

	if !args.Present() {
		fmt.Fprintln(ctx.App.Writer, "Filename required")
		fmt.Fprintln(ctx.App.Writer)

		cli.ShowCommandHelp(ctx, "dump")
		return
	}

	rrdFile, err := cdata.OpenRrdRawFile(args.First(), true)
	if err != nil {
		fmt.Fprintln(ctx.App.Writer, err)
		return
	}
	defer rrdFile.Close()
}
