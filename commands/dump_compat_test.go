package commands

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
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
			"--start", "now",
			"--step", "1s",
			"DS:watts:GAUGE:5m:0:100000",
			"RRA:AVERAGE:0.5:5s:60m")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		So(cmd.Run(), ShouldBeNil)

		Convey("Then dump produces the same result", func() {
			cmd := exec.Command(rrdtool, "dump", rrdFileName)
			stdout, err := cmd.StdoutPipe()

			So(err, ShouldBeNil)

			rrdDecoder := xml.NewDecoder(stdout)
			cmd.Start()

			for {
				token, err := rrdDecoder.Token()
				if err == io.EOF {
					break
				}
				Printf("%#v", token)
			}
		})
	})
}
