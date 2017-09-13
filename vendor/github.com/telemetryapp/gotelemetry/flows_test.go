package gotelemetry

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Example_flow() {
	c, _ := NewCredentials("test-api-token")

	g := Gauge{
		Value: 123,
	}

	f := NewFlow("test_gauge", &g)

	err := f.Publish(c)

	if err != nil {
		panic(fmt.Sprintf("Something went wrong: %v\n", err))
	}
}

func ExampleFlow() {
	c, _ := NewCredentials("test-api-token")

	g := Gauge{
		Value: 123,
	}

	f := NewFlow("test_gauge", &g)

	err := f.Publish(c)

	if err != nil {
		panic(fmt.Sprintf("Something went wrong: %v\n", err))
	}
}

func TestFlows(t *testing.T) {
	Convey("Flow objects", t, func() {

		Convey("Should support instantiation", func() {
			g := Gauge{
				Value: 123,
			}

			f := NewFlow("test_gauge", &g)

			So(f, ShouldNotBeNil)
			So(f.Tag, ShouldEqual, "test_gauge")
			So(*(f.Data.(*Gauge)), ShouldResemble, g)
		})

		Convey("Should support posting", func() {
			c, err := NewCredentials("test-api-token")

			So(err, ShouldBeNil)

			g := Gauge{
				Value: 123,
			}

			f := NewFlow("test-flow-gauge", &g)

			err = f.Publish(c)

			So(err, ShouldBeNil)
		})

		Convey("Should properly report errors", func() {
			c, err := NewCredentials("test-api-token")

			So(err, ShouldBeNil)

			g := Gauge{}

			f := NewFlow("test-flow-test", &g)

			err = f.Publish(c)

			So(err, ShouldNotBeNil)

			e, ok := err.(*Error)

			So(ok, ShouldBeTrue)
			So(e.StatusCode, ShouldEqual, 400)
		})

		Convey("Should support retrieving data by variant type and moddifying it", func() {
			g := Gauge{
				Value: 123,
			}

			f := NewFlow("test_gauge", &g)

			gg, _ := f.GaugeData()

			So(*gg, ShouldResemble, g)

			gg.Value = 124

			So(f.Data.(*Gauge).Value, ShouldEqual, 124)
		})

		Convey("Should require pointers to variant structures to be passed on init", func() {
			f := func() {
				g := Gauge{
					Value: 123,
				}

				NewFlow("test_gauge", g)
			}

			So(f, ShouldPanic)
		})
	})
}
