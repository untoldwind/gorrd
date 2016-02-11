package cdata

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnival(t *testing.T) {
	Convey("Unival from uint64", t, func() {
		val := unival(123456)

		So(val.AsLong(), ShouldEqual, 123456)
		So(val.AsUnsignedLong(), ShouldEqual, 123456)
	})

	Convey("Unival from float64", t, func() {
		val := univalForDouble(1234.678)

		So(val.AsDouble(), ShouldEqual, 1234.678)
	})
}
