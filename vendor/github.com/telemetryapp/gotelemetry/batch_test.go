package gotelemetry

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Example_batch() {
	// Note that the test API token cannot be used for
	// batch submission.
	c, err := NewCredentials(getTestKey())

	So(err, ShouldBeNil)

	g := Gauge{Value: 10}
	fg := NewFlow("gauge_1", &g)

	v := Value{Value: 101.23}
	fv := NewFlow("gauge_2", &v)

	b := Batch{}

	b.SetFlow(fg)
	b.SetFlow(fv)

	err = b.Publish(c)

	if err != nil {
		panic(fmt.Sprintf("Something went wrong: %v\n", err))
	}
}

func ExampleBatch() {
	// Note that the test API token cannot be used for
	// batch submission.
	c, err := NewCredentials(getTestKey())

	So(err, ShouldBeNil)

	g := Gauge{Value: 10}
	fg := NewFlow("gauge_1", &g)

	v := Value{Value: 101.23}
	fv := NewFlow("gauge_2", &v)

	b := Batch{}

	b.SetFlow(fg)
	b.SetFlow(fv)

	err = b.Publish(c)

	if err != nil {
		panic(fmt.Sprintf("Something went wrong: %v\n", err))
	}
}

func TestBatchSubmissions(t *testing.T) {
	Convey("The batch submission system", t, func() {

		Convey("Should allow creating batches", func() {
			g := Gauge{Value: 10}
			fg := NewFlow("test-flow-gauge", &g)

			v := Value{Value: 10}
			fv := NewFlow("test-flow-value", &v)

			b := Batch{}

			b.SetFlow(fg)
			b.SetFlow(fv)

			So(len(b), ShouldEqual, 2)
		})

		Convey("Should allow retrieving flows", func() {
			g := Gauge{Value: 10}
			fg := NewFlow("test-flow-gauge", &g)

			v := Value{Value: 10}
			fv := NewFlow("test-flow-value", &v)

			b := Batch{}

			b.SetFlow(fg)
			b.SetFlow(fv)

			So(len(b), ShouldEqual, 2)

			x, ok := b.Flow("test-flow-gauge")

			So(ok, ShouldBeTrue)
			So(x, ShouldResemble, fg)

			x, ok = b.Flow("test-flow-value")

			So(ok, ShouldBeTrue)
			So(x, ShouldResemble, fv)

			x, ok = b.Flow("BLARG")

			So(ok, ShouldBeFalse)
			So(x, ShouldBeNil)
		})

		Convey("Should post data to the server and return an error", func() {
			c, err := NewCredentials("test-api-token")

			So(err, ShouldBeNil)

			g := Gauge{Value: 10}
			fg := NewFlow("test-flow-gauge", &g)

			v := Value{Value: 10}
			fv := NewFlow("test-flow-value", &v)

			b := Batch{}

			b.SetFlow(fg)
			b.SetFlow(fv)

			err = b.Publish(c)

			So(err, ShouldNotBeNil)
			So(err.(*Error).StatusCode, ShouldEqual, 401)
		})
	})

}
