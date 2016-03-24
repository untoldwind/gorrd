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

		Convey("When values are added with gaps", func() {
			rrdFileNameCopy := filepath.Join(tempDir, fmt.Sprintf("comp_update1-copy-%d.rrd", time.Now().UnixNano()))
			//			defer os.Remove(rrdFileNameCopy)
			fmt.Println(rrdFileNameCopy)

			timestamp := int64(start + 400)
			for i := 0; i < 20; i++ {
				copyFile(rrdFileName, rrdFileNameCopy)
				So(rrdtool.update(
					rrdFileName,
					fmt.Sprintf("%d:%d", timestamp, i*100+5),
				), ShouldBeNil)

				err := runUpdateCommand(rrdFileNameCopy, time.Unix(timestamp, 0), strconv.Itoa(i*100+5))

				So(err, ShouldBeNil)
				So(rrdFileNameCopy, shouldHaveSameContentAs, rrdFileName)

				if i%4 < 3 {
					timestamp += 200
				} else {
					timestamp += 400
				}
			}
		})
	})

	Convey("Given minimal rrdfile with 5m step", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update2-%d.rrd", time.Now().UnixNano()))
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
			rrdFileNameCopy := filepath.Join(tempDir, fmt.Sprintf("comp_update2-copy-%d.rrd", time.Now().UnixNano()))
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
