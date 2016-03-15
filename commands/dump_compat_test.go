package commands_test

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/codegangsta/cli"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/untoldwind/gorrd/commands"
)

func TestDumpCompatibility(t *testing.T) {
	rrdtool, err := exec.LookPath("rrdtool")

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}
	Convey("Given minimal rrdfile with 1s step", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update-%s.rrd", time.Now().String()))
		defer os.Remove(rrdFileName)

		cmd := exec.Command(rrdtool,
			"create",
			rrdFileName,
			"--start", "1455218381",
			"--step", "1",
			"DS:watts:GAUGE:300:0:100000",
			"RRA:AVERAGE:0.5:5:3600")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		So(cmd.Run(), ShouldBeNil)

		Convey("Then dump produces the same result", func() {
			cmd := exec.Command(rrdtool, "dump", rrdFileName)
			stdout, err := cmd.StdoutPipe()

			So(err, ShouldBeNil)

			cmd.Start()
			expectedResult, err := flattenXml(stdout)

			So(err, ShouldBeNil)

			pipeReader, pipeWriter := io.Pipe()
			go func() {
				flags := flag.NewFlagSet("gorrd", flag.ContinueOnError)
				flags.Parse([]string{rrdFileName})
				ctx := cli.NewContext(&cli.App{
					Writer: pipeWriter,
				}, flags, nil)
				commands.DumpCommand.Action(ctx)
				pipeWriter.Close()
			}()

			actualResult, err := flattenXml(pipeReader)

			So(err, ShouldBeNil)

			So(actualResult, ShouldResemble, expectedResult)
		})
	})
}
