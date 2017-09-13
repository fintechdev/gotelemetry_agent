package gotelemetry

import (
	"crypto/rand"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBoards(t *testing.T) {
	credentials, _ := NewCredentials(getTestKey())

	p := make([]byte, 10)
	rand.Read(p)

	name := fmt.Sprintf("board-%x", p)

	Convey("Boards", t, func() {
		Convey("Should properly create, retrieve, and delete boards", func() {
			b, err := NewBoard(credentials, name, "dark", true, "HDTV")

			So(err, ShouldBeNil)
			So(b, ShouldNotBeNil)
			So(b.Id, ShouldNotBeNil)
			So(b.ChannelIds, ShouldNotBeEmpty)

			b1, err := GetBoard(credentials, b.Id)

			So(err, ShouldBeNil)
			So(b1, ShouldNotBeNil)
			So(b1.Id, ShouldEqual, b.Id)

			err = b.Delete()
			So(err, ShouldBeNil)
		})

		Convey("Should fail to retrieve and delete a invalid board", func() {
			newBoard, err := NewBoard(credentials, name, "dark", true, "HDTV")

			newBoard.Id += "123"

			getBoard, err := GetBoard(credentials, newBoard.Id)
			So(err, ShouldNotBeNil)
			So(getBoard, ShouldBeNil)

			err = newBoard.Delete()
			So(err, ShouldNotBeNil)

		})
	})

}
