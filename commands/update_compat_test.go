package commands_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/codegangsta/cli"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/untoldwind/gorrd/commands"
)

func TestUpdateCompatibility(t *testing.T) {
	rrdtool, err := findRrdTool()

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	runUpdateCommand := func(args ...string) {
		flags := flag.NewFlagSet("gorrd", flag.ContinueOnError)
		flags.Parse(args)
		ctx := cli.NewContext(&cli.App{
			Writer: os.Stdout,
		}, flags, nil)
		commands.UpdateCommand.Action(ctx)
	}

	Convey("Given minimal rrdfile with 5m step", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update-%d.rrd", time.Now().UnixNano()))
		defer os.Remove(rrdFileName)

		start := 1455218381
		So(rrdtool.create(
			rrdFileName,
			strconv.Itoa(start),
			"300",
			"DS:watts:GAUGE:300:0:100000",
			"RRA:AVERAGE:0.5:12:24",
		), ShouldBeNil)

		Convey("When values are added within stepsize", func() {
			rrdFileNameCopy := filepath.Join(tempDir, fmt.Sprintf("comp_update-copy-%d.rrd", time.Now().UnixNano()))
			defer os.Remove(rrdFileNameCopy)

			for i := 1; i < 5; i++ {
				copyFile(rrdFileName, rrdFileNameCopy)
				So(rrdtool.update(
					rrdFileName,
					fmt.Sprintf("%d:%d", i+start, i*100+5),
				), ShouldBeNil)

				runUpdateCommand(rrdFileNameCopy, fmt.Sprintf("%d:%d", i+start, i*100+5))

				So(rrdFileNameCopy, shouldHaveSameContentAs, rrdFileName)
			}
		})
	})
}
