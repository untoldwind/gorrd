package rrd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/untoldwind/gorrd/rrd/cdata"
)

func TestUpdateCompatibility(t *testing.T) {
	rrdtool, err := findRrdTool()

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	runUpdateCommand := func(filename string, timestamp time.Time, values ...string) error {
		rrd, err := cdata.OpenRrdRawFile(filename, false)
		if err != nil {
			return err
		}
		defer rrd.Close()

		return rrd.Update(timestamp, values)
	}

	Convey("Given minimal rrdfile with 5m step and 1 pdp per cdp", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update1-%d.rrd", time.Now().UnixNano()))
		defer os.Remove(rrdFileName)

		start := 1455218381
		So(rrdtool.create(
			rrdFileName,
			strconv.Itoa(start),
			"300",
			"DS:watts:GAUGE:300:0:100000",
			"RRA:AVERAGE:0.5:1:24",
		), ShouldBeNil)

		Convey("When values are added in stepsize", func() {
			rrdFileNameCopy := filepath.Join(tempDir, fmt.Sprintf("comp_update1-copy-%d.rrd", time.Now().UnixNano()))
			defer os.Remove(rrdFileNameCopy)

			for i := 1; i < 10; i++ {
				copyFile(rrdFileName, rrdFileNameCopy)
				So(rrdtool.update(
					rrdFileName,
					fmt.Sprintf("%d:%d", 300*i+start, i*100+5),
				), ShouldBeNil)

				err := runUpdateCommand(rrdFileNameCopy, time.Unix(int64(300*i+start), 0), strconv.Itoa(i*100+5))

				So(err, ShouldBeNil)
				So(rrdFileNameCopy, shouldHaveSameContentAs, rrdFileName)
			}
		})

		Convey("When values are added above stepsize", func() {
			rrdFileNameCopy := filepath.Join(tempDir, fmt.Sprintf("comp_update1-copy-%d.rrd", time.Now().UnixNano()))
			defer os.Remove(rrdFileNameCopy)

			for i := 1; i < 10; i++ {
				copyFile(rrdFileName, rrdFileNameCopy)
				So(rrdtool.update(
					rrdFileName,
					fmt.Sprintf("%d:%d", 400*i+start, i*100+5),
				), ShouldBeNil)

				err := runUpdateCommand(rrdFileNameCopy, time.Unix(int64(400*i+start), 0), strconv.Itoa(i*100+5))

				So(err, ShouldBeNil)
				So(rrdFileNameCopy, shouldHaveSameContentAs, rrdFileName)
			}
		})
	})

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

				err := runUpdateCommand(rrdFileNameCopy, time.Unix(int64(i+start), 0), strconv.Itoa(i*100+5))

				So(err, ShouldBeNil)
				So(rrdFileNameCopy, shouldHaveSameContentAs, rrdFileName)
			}
		})
	})
}
