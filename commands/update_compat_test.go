package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
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
			"--step", "5s",
			"DS:watts:GAUGE:5m:0:100000",
			"RRA:AVERAGE:0.5:5s:60m")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		So(cmd.Run(), ShouldBeNil)
	})
}
