package utils

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetHumanSize(t *testing.T) {
	Convey("Test GetHumanSize", t, func() {
		Convey("T", func() {
			var v = uint64(25 * 1024 * 1024 * 1024 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "25.0T")
		})
		Convey("G", func() {
			var v = uint64(33880356 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "32.3G")
		})
		Convey("M", func() {
			var v = uint64(5812 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "5.7M")
		})
		Convey("K", func() {
			var v = uint64(123 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "123.0K")
		})
		Convey("B", func() {
			var v = uint64(786)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "786B")
		})
		Convey("0", func() {
			var v = uint64(0)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "0B")
		})
	})
}

func Test_AddOneDay(t *testing.T) {
	Convey("Test AddOneDay", t, func() {
		Convey("T", func() {
			var v = "2016-02-29"
			ret := AddOneDay(v)
			So(ret, ShouldEqual, "2016-03-01")
		})

	})
}
