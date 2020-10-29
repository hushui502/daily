package goconvey

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func ShouldCase(actual interface{}, expected ...interface{}) string {
	if actual == expected[0] {
		return ""
	} else {
		return "n"
	}
}

func TestAdd(t *testing.T) {
	Convey("相加", t, func() {
		So(Add(1, 3), ShouldCase, 4)
	})
}
