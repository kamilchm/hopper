package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"io/ioutil"
	"log"
)

func TestConfig(t *testing.T) {
	Convey("Given config with echo docker hop", t, func() {
		hopDef := `
echo:
- docker:
    image: ubuntu
    command: echo
`
		configFile := prepareConfig(hopDef)

		Convey("When hops are parsed", func() {
			hops, _ := LoadHops(configFile)
			echoHop, present := hops["echo"]

			Convey("There should be echo hop", func() {
				So(present, ShouldBeTrue)
			})

			Convey("The hop should be docker type", func() {
				So(echoHop[0], ShouldHaveSameTypeAs, Docker{})
			})
		})
	})
}

func prepareConfig(content string) string {
	tempFile, err := ioutil.TempFile("", "hopper-test")
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()
	tempFile.WriteString(content)

	return tempFile.Name()
}
