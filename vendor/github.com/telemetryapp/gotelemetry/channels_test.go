package gotelemetry

import (
	"crypto/rand"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestChannels(t *testing.T) {
	credentials, _ := NewCredentials(getTestKey())

	p := make([]byte, 10)
	rand.Read(p)

	name := fmt.Sprintf("board-%x", p)

	Convey("Board export", t, func() {
		Convey("Should properly export a board", func() {
			b, err := NewBoard(credentials, name, "dark", true, "HDTV")
			_, err = NewWidget(credentials, b, "value", 1, 2, 3, 4, 0, "normal")
			_, err = NewWidget(credentials, b, "value", 1, 2, 3, 4, 0, "normal")

			b1, err := GetBoard(credentials, b.Id)

			exported, err := b1.Export()

			So(err, ShouldBeNil)
			So(exported, ShouldNotBeNil)

			b1.Delete()

			b3, err := ImportBoard(credentials, name+"-clone", "testt", exported)

			So(err, ShouldBeNil)
			So(b3, ShouldNotBeNil)

			b4, err := GetBoard(credentials, b3.Id)

			So(len(b3.Widgets), ShouldEqual, len(b4.Widgets))

			b3.Delete()
		})

		Convey("Should return, rather than reimport, a board with a name that already exists", func() {

			newBoard, _ := NewBoard(credentials, name, "dark", true, "HDTV")
			_, _ = NewWidget(credentials, newBoard, "value", 1, 2, 3, 4, 0, "normal")
			_, _ = NewWidget(credentials, newBoard, "value", 1, 2, 3, 4, 0, "normal")

			getBoard, _ := GetBoard(credentials, newBoard.Id)

			exportedBoard, _ := getBoard.Export()

			importedBoard, err := ImportBoard(credentials, "Fail Import", "fail_import", exportedBoard)
			importedBoardWithSameName, err := ImportBoard(credentials, "Fail Import", "fail_import", exportedBoard)

			So(importedBoardWithSameName, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(importedBoardWithSameName.Id, ShouldEqual, importedBoard.Id)

			// Clean up
			importedBoard.Delete()

		})
	})

}
