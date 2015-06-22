package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"io/ioutil"
	"os"
)

func TestHopper(t *testing.T) {
	Convey("Given empty dir", t, func() {
		testDir := createTestDir()

		Convey("When we go there", func() {
			os.Chdir(testDir)

			Convey("We shouldn't be in local mode", func() {
				So(inLocalMode(), ShouldBeFalse)
			})
		})
	})

	Convey("Given dir with hop.yaml", t, func() {
		testDir := createTestDir()
		ioutil.WriteFile(testDir+"/hop.yaml", []byte(""), 0644)

		Convey("When we go there", func() {
			os.Chdir(testDir)

			Convey("We should be in local mode", func() {
				So(inLocalMode(), ShouldBeTrue)
			})
		})
	})
}
