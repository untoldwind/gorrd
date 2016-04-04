package commands_test

import (
	"bytes"
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

func TestDumpCompatibility(t *testing.T) {
	rrdtool, err := findRrdTool()

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	runDumpCommand := func(rrdFileName string) (map[string]interface{}, error) {
		output := bytes.NewBufferString("")
		flags := flag.NewFlagSet("gorrd", flag.ContinueOnError)
		flags.Parse([]string{rrdFileName})
		ctx := cli.NewContext(&cli.App{
			Writer: output,
		}, flags, nil)
		commands.DumpCommand.Action(ctx)

		return flattenXml(output.String())
	}

	Convey("Given minimal rrdfile with 5m step", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update1-%d.rrd", time.Now().UnixNano()))
		defer os.Remove(rrdFileName)

		start := 1455218381
		err := rrdtool.create(rrdFileName,
			strconv.Itoa(start),
			"1",
			"DS:watts:GAUGE:300:0:100000",
			"RRA:AVERAGE:0.5:1:100",
		)

		So(err, ShouldBeNil)

		expectedResult, err := rrdtool.dump(rrdFileName)

		So(err, ShouldBeNil)

		actualResult, err := runDumpCommand(rrdFileName)

		So(err, ShouldBeNil)
		So(expectedResult, ShouldResemble, actualResult)
	})
}
