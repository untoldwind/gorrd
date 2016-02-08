package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/untoldwind/gorrd/rrd/cdata"
	"github.com/untoldwind/gorrd/rrd/dump"
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
	filename, err := getFilenameArg(ctx)
	if err != nil {
		showError(ctx, err)
		return
	}

	rrd, err := cdata.OpenRrdRawFile(filename, true)
	if err != nil {
		fmt.Fprintln(ctx.App.Writer, err)
		return
	}
	defer rrd.Close()

	xmlDumper, err := dump.NewXmlOutput(ctx.App.Writer, true)
	if err != nil {
		fmt.Fprintln(ctx.App.Writer, err)
		return
	}
	if err := rrd.DumpTo(xmlDumper); err != nil {
		fmt.Fprintln(ctx.App.Writer, err)
		return
	}
}
