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
