package rrd

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUtils(t *testing.T) {
	Convey("rrdIsUnsignedInt", t, func() {
		So(rrdIsUnsignedInt("0"), ShouldBeTrue)
		So(rrdIsUnsignedInt("123456789"), ShouldBeTrue)
		So(rrdIsUnsignedInt("-1"), ShouldBeFalse)
		So(rrdIsUnsignedInt("1abc"), ShouldBeFalse)
	})

	Convey("rrdIsSignedInt", t, func() {
		So(rrdIsSignedInt("0"), ShouldBeTrue)
		So(rrdIsSignedInt("123456789"), ShouldBeTrue)
		So(rrdIsSignedInt("-1"), ShouldBeTrue)
		So(rrdIsSignedInt("-123456789"), ShouldBeTrue)
		So(rrdIsSignedInt("1abc"), ShouldBeFalse)
	})

	Convey("rrdDiff", t, func() {
		So(rrdDiff("0", "0"), ShouldEqual, 0)
		So(rrdDiff("123", "23"), ShouldEqual, 100)
		So(rrdDiff("23", "123"), ShouldEqual, -100)
		So(rrdDiff("1234567890", "987654321"), ShouldEqual, 246913569)
		So(rrdDiff("987654321", "1234567890"), ShouldEqual, -246913569)
		So(rrdDiff("-1234567890", "987654321"), ShouldEqual, -246913569)
		So(rrdDiff("987654321", "1234567890"), ShouldEqual, -246913569)
	})
}