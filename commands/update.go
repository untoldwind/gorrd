package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
	"github.com/untoldwind/gorrd/rrd/cdata"
)

// Command line create command
var UpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update values to an rrd rrd file",
	Flags:     []cli.Flag{},
	ArgsUsage: "file.rrd timestamp:value[:value]",
	Action:    updateCommand,
}

func updateCommand(ctx *cli.Context) {
	filename, err := getFilenameArg(ctx)
	if err != nil {
		showError(ctx, err)
		return
	}

	rrd, err := cdata.OpenRrdRawFile(filename, true)
	if err != nil {
		showError(ctx, err)
		return
	}
	defer rrd.Close()

	timestamp, values, err := parseUpdateArg(ctx, rrd)
	if err != nil {
		showError(ctx, err)
		return
	}

	if err := rrd.Update(timestamp, values); err != nil {
		showError(ctx, err)
		return
	}
}

func parseUpdateArg(ctx *cli.Context, rrd *rrd.Rrd) (time.Time, []string, error) {
	var timestamp time.Time
	if len(ctx.Args()) < 2 {
		return timestamp, nil, errors.Errorf("Update argument required")
	}
	parts := strings.Split(ctx.Args().Get(1), ":")

	if len(parts) != len(rrd.Datasources)+1 {
		return time.Time{}, nil, errors.Errorf("expected %d data source readings (got %d)", len(rrd.Datasources), len(parts)-1)
	}
	if parts[0] == "N" {
		timestamp = time.Now()
	} else {
		unixTime, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return timestamp, nil, errors.Errorf("Update argument required")
		}
		timestamp = time.Unix(unixTime, 0)
	}
	return timestamp, parts[1:], nil
}
