package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"io/ioutil"
	"log"
)

func TestDockerConfig(t *testing.T) {
	Convey("Given config with two docker hops", t, func() {
		hopDef := `---
echo:
- docker:
    image: ubuntu
    command: echo
cat:
- docker:
    image: alpine
    command: cat
    permissions:
        cwd: yes
`
		configFile := prepareConfig(hopDef)

		Convey("When hops are parsed", func() {
			hops, _ := LoadHops(configFile)

			echoHop, present := hops["echo"]
			Convey("There should be echo hop", func() {
				So(present, ShouldBeTrue)
			})

			catHop, present := hops["cat"]
			Convey("There should be cat hop", func() {
				So(present, ShouldBeTrue)
			})

			Convey("Both hops should have docker definition", func() {
				So(echoHop[0], ShouldResemble, &Docker{
					Image:   "ubuntu",
					Command: "echo",
				})
				So(catHop[0], ShouldResemble, &Docker{
					Image:       "alpine",
					Command:     "cat",
					Permissions: permissions{true},
				})
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
