package rrd_test

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/untoldwind/gorrd/rrd"
)

func TestDatasourceCounter(t *testing.T) {
	Convey("Given unbounded COUNTER datasource", t, func() {
		datasource := &rrd.DatasourceCounter{
			DatasourceAbstract: rrd.DatasourceAbstract{
				Name:      "test",
				Heartbeat: 300,
				Min:       math.NaN(),
				Max:       math.NaN(),
			},
		}

		Convey("When datasource has unknown last value", func() {
			datasource.LastValue = "U"

			Convey("When value is updated within heartbeat", func() {
				newPdp, err := datasource.CalculatePdpPrep("10000", 100)

				So(err, ShouldBeNil)
				So(math.IsNaN(newPdp), ShouldBeTrue)
				So(datasource.LastValue, ShouldEqual, "10000")
			})

			Convey("When value is updated beyond heartbeat", func() {
				newPdp, err := datasource.CalculatePdpPrep("10000", 1000)

				So(err, ShouldBeNil)
				So(math.IsNaN(newPdp), ShouldBeTrue)
				So(datasource.LastValue, ShouldEqual, "10000")
			})
		})

		Convey("When datasource has last value 123456", func() {
			datasource.LastValue = "123456"

			Convey("When value 234567 is updated within heartbeat", func() {
				newPdp, err := datasource.CalculatePdpPrep("234567", 100)

				So(err, ShouldBeNil)
				So(newPdp, ShouldEqual, 111111)
				So(datasource.LastValue, ShouldEqual, "234567")
			})

			Convey("When value 234567 is updated beyond heartbeat", func() {
				newPdp, err := datasource.CalculatePdpPrep("234567", 1000)

				So(err, ShouldBeNil)
				So(math.IsNaN(newPdp), ShouldBeTrue)
				So(datasource.LastValue, ShouldEqual, "234567")
			})
		})
	})
}
