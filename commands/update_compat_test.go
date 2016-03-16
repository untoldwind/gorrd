package commands_test

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/codegangsta/cli"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/untoldwind/gorrd/commands"
)

func TestUpdateCompatibility(t *testing.T) {
	rrdtool, err := exec.LookPath("rrdtool")

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	Convey("Given minimal rrdfile with 5s step", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update-%s.rrd", time.Now().String()))
		defer os.Remove(rrdFileName)

		cmd := exec.Command(rrdtool,
			"create",
			rrdFileName,
			"--start", "1455218381",
			"--step", "300",
			"DS:watts:GAUGE:300:0:100000",
			"RRA:AVERAGE:0.5:12:24")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		So(cmd.Run(), ShouldBeNil)

		Convey("When values are added within stepsize", func() {
			rrdFileNameCopy := filepath.Join(tempDir, fmt.Sprintf("comp_update-copy-%s.rrd", time.Now().String()))
			defer os.Remove(rrdFileNameCopy)

			for i := 1; i < 5; i++ {
				copyFile(rrdFileName, rrdFileNameCopy)
				cmd := exec.Command(rrdtool,
					"update",
					rrdFileName,
					fmt.Sprintf("%d:%d", i+1455218381, i*100+5))
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				So(cmd.Run(), ShouldBeNil)

				flags := flag.NewFlagSet("gorrd", flag.ContinueOnError)
				flags.Parse([]string{rrdFileNameCopy, fmt.Sprintf("%d:%d", i+1455218381, i*100+5)})
				ctx := cli.NewContext(&cli.App{
					Writer: os.Stdout,
				}, flags, nil)
				commands.UpdateCommand.Action(ctx)

				So(rrdFileNameCopy, shouldHaveSameContentAs, rrdFileName)
			}
		})
	})
}
