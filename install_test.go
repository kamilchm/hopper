package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"fmt"
	"os"
	"syscall"
	"time"
)

func TestInstall(t *testing.T) {
	Convey("Given empty dir", t, func() {
		wrk := workspace{Hops: nil, BinDir: createTestDir(),
			HopperPath: "/bin/hopper"}

		Convey("When installing hop", func() {
			So(wrk.installHop("some-hop", false), ShouldBeNil)

			Convey("It should link to hopper", func() {
				So(wrk.BinDir+"/some-hop", shoouldLinkTo, "/bin/hopper")
			})
		})

	})

	Convey("Given dir with echo hop", t, func() {
		wrk := workspace{Hops: nil, BinDir: createTestDir(),
			HopperPath: os.Args[0]}
		So(wrk.installHop("echo", false), ShouldBeNil)
		echoBefore := hopTime("echo", wrk.BinDir)
		time.Sleep(5 * time.Millisecond)

		Convey("When installing echo and cat hop", func() {
			So(wrk.installHop("echo", false), ShouldBeNil)
			So(wrk.installHop("cat", false), ShouldBeNil)

			Convey("cat should link to hopper", func() {
				So(wrk.BinDir+"/cat", shoouldLinkTo, os.Args[0])
			})

			Convey("echo should remain unchanged", func() {
				So(hopTime("echo", wrk.BinDir), ShouldResemble, echoBefore)
			})

			Convey("It could be installed with the --force flag", func() {
				So(wrk.installHop("echo", true), ShouldBeNil)
				Convey("and echo shouldn't be updated", func() {
					So(hopTime("echo", wrk.BinDir), ShouldResemble, echoBefore)
				})
			})

		})

	})

	Convey("Given dir with some file named echo", t, func() {
		wrk := workspace{Hops: nil, BinDir: createTestDir(),
			HopperPath: "/bin/bash"}
		So(wrk.installHop("echo", false), ShouldBeNil)
		echoBefore := hopTime("echo", wrk.BinDir)
		time.Sleep(5 * time.Millisecond)

		wrk.HopperPath = os.Args[0]
		Convey("When installing echo hop", func() {
			So(wrk.installHop("echo", false), ShouldNotBeNil)

			Convey("echo file should remain unchanged", func() {
				So(hopTime("echo", wrk.BinDir), ShouldResemble, echoBefore)
			})

			Convey("but when installing with the --force flag", func() {
				So(wrk.installHop("echo", true), ShouldBeNil)

				Convey("echo should be replaced", func() {
					So(hopTime("echo", wrk.BinDir), ShouldHappenAfter,
						echoBefore)
				})
			})

		})

	})
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
	}
	return fmt.Sprintf("Expected target: %v, but it links to %v",
		expected, linkTarget)
}
