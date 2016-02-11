package cdata

import (
	"encoding/binary"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestRrdRawHeader(t *testing.T) {
	rrdtool, err := exec.LookPath("rrdtool")

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	tempDir := os.TempDir()
	rrdFileName := filepath.Join(tempDir, fmt.Sprintf("minimal%d.rrd", time.Now().Unix()))
	defer os.Remove(rrdFileName)

	Convey("Given minimal rrdfile", t, func() {
		cmd := exec.Command(rrdtool,
			"create",
			rrdFileName,
			"--start", "1455211143",
			"--step", "1s",
			"DS:watts:GAUGE:5m:0:100000",
			"RRA:AVERAGE:0.5:1s:5m",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		So(cmd.Run(), ShouldBeNil)

		dataFile, err := OpenCDataFile(rrdFileName, true, binary.LittleEndian, 8, 8)

		So(err, ShouldBeNil)
		So(dataFile, ShouldNotBeNil)
		defer dataFile.Close()

		rrdFile := &RrdRawFile{
			dataFile: dataFile,
		}
		reader := dataFile.Reader(0)
		So(rrdFile.readHeaders(reader), ShouldBeNil)

		So(rrdFile.header.datasourceCount, ShouldEqual, 1)
		So(rrdFile.header.rraCount, ShouldEqual, 1)
		So(rrdFile.header.pdpStep, ShouldEqual, 1)
		So(rrdFile.datasourceDefs, ShouldHaveLength, 1)
		So(rrdFile.datasourceDefs[0].name, ShouldEqual, "watts")
		So(rrdFile.datasourceDefs[0].dataSourceType, ShouldEqual, "GAUGE")
		So(rrdFile.rraDefs, ShouldHaveLength, 1)
		So(rrdFile.rraDefs[0].rraType, ShouldEqual, "AVERAGE")
		So(rrdFile.rraDefs[0].pdpPerRow, ShouldEqual, 1)
		So(rrdFile.rraDefs[0].rowCount, ShouldEqual, 300)
		So(rrdFile.lastUpdate, ShouldResemble, time.Unix(1455211143, 0))
	})
}
