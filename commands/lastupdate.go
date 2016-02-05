package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/untoldwind/gorrd/rrd/cdata"
)

// Command line create command
var LastUpdateCommand = cli.Command{
	Name:      "lastupdate",
	Usage:     "Show the last update to a rrd file",
	Flags:     []cli.Flag{},
	ArgsUsage: "file.rrd",
	Action:    lastUpdate,
}

func lastUpdate(ctx *cli.Context) {
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

	for _, datasource := range rrd.Datasources {
		fmt.Fprintf(ctx.App.Writer, " %s", datasource.GetName())
	}
	fmt.Fprint(ctx.App.Writer, "\n\n")

	fmt.Fprintf(ctx.App.Writer, "%10d:", rrd.LastUpdate.Unix())
	for _, datasource := range rrd.Datasources {
		fmt.Fprintf(ctx.App.Writer, " %s", datasource.GetLastValue())
	}
	fmt.Fprint(ctx.App.Writer, "\n")
}
