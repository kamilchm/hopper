package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"
)

func TestInstall(t *testing.T) {
	Convey("Given empty dir", t, func() {
		testDir := createTestDir()

		Convey("When installing hop", func() {
			So(installHop("some-hop", testDir, "/bin/hopper", false),
				ShouldBeNil)

			Convey("It should link to hopper", func() {
				So(testDir+"/some-hop", shoouldLinkTo, "/bin/hopper")
			})
		})

	})

	Convey("Given dir with echo hop", t, func() {
		testDir := createTestDir()
		So(installHop("echo", testDir, os.Args[0], false), ShouldBeNil)
		echoBefore := hopTime("echo", testDir)
		time.Sleep(5 * time.Millisecond)

		Convey("When installing echo and cat hop", func() {
			So(installHop("echo", testDir, os.Args[0], false), ShouldBeNil)
			So(installHop("cat", testDir, os.Args[0], false), ShouldBeNil)

			Convey("cat should link to hopper", func() {
				So(testDir+"/cat", shoouldLinkTo, os.Args[0])
			})

			Convey("echo should remain unchanged", func() {
				So(hopTime("echo", testDir), ShouldResemble, echoBefore)
			})

			Convey("It could be installed with the --force flag", func() {
				So(installHop("echo", testDir, os.Args[0], true),
					ShouldBeNil)
				Convey("and echo shouldn't be updated", func() {
					So(hopTime("echo", testDir), ShouldResemble, echoBefore)
				})
			})

		})

	})

	Convey("Given dir with some file named echo", t, func() {
		testDir := createTestDir()
		So(installHop("echo", testDir, "/bin/bash", false), ShouldBeNil)
		echoBefore := hopTime("echo", testDir)
		time.Sleep(5 * time.Millisecond)

		Convey("When installing echo hop", func() {
			So(installHop("echo", testDir, os.Args[0], false),
				ShouldNotBeNil)

			Convey("echo file should remain unchanged", func() {
				So(hopTime("echo", testDir), ShouldResemble, echoBefore)
			})

			Convey("but when installing with the --force flag", func() {
				So(installHop("echo", testDir, os.Args[0], true),
					ShouldBeNil)

				Convey("echo should be replaced", func() {
					So(hopTime("echo", testDir), ShouldHappenAfter,
						echoBefore)
				})
			})

		})

	})
}

func createTestDir() string {
	emptyDir, err := ioutil.TempDir("", "hopper-test")
	if err != nil {
		log.Fatal(err)
	}
	return emptyDir
}

func hopTime(name, dir string) time.Time {
	hopInfo, err := os.Lstat(dir + "/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return timespecToTime(hopInfo.Sys().(*syscall.Stat_t).Ctim)
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func shoouldLinkTo(actual interface{}, expected ...interface{}) string {
	linkTarget, err := os.Readlink(actual.(string))
	if err != nil {
		log.Fatal(err)
	}

	if linkTarget == expected[0].(string) {
		return ""
	} else {
		return fmt.Sprintf("Expected target: %v, but it links to %v",
			expected, linkTarget)
	}
}
