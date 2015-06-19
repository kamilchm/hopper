package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSearch(t *testing.T) {
	Convey("Given no hops", t, func() {
		noHops := hops{}

		Convey("When searching exact hop", func() {
			_, err := noHops.searchHop("some-hop")

			Convey("It should fail with NotFound error", func() {
				So(err, ShouldEqual, HopNotFound)
			})
		})

	})

	Convey("Given hops: hop1, hop3, some", t, func() {
		defPerms := permissions{}
		theHops := hops{
			"hop1": []Hop{&Docker{"a", "hop1", defPerms}},
			"hop3": []Hop{&Docker{"a", "hop3", defPerms}},
			"some": []Hop{&Docker{"a", "some", defPerms}},
		}

		Convey("When search for hop1", func() {
			found, _ := theHops.searchHop("hop1")

			Convey("It should return only hop1", func() {
				So(found, ShouldResemble, &hops{
					"hop1": []Hop{&Docker{"a", "hop1", defPerms}},
				})
			})
		})

		Convey("When search for hop2", func() {
			_, err := theHops.searchHop("hop2")

			Convey("It should fail with NotFound error", func() {
				So(err, ShouldEqual, HopNotFound)
			})
		})

		Convey("When search for hop*", func() {
			found, _ := theHops.searchHop("hop*")

			Convey("It should return hop1 and hop3", func() {
				So(found, ShouldResemble, &hops{
					"hop1": []Hop{&Docker{"a", "hop1", defPerms}},
					"hop3": []Hop{&Docker{"a", "hop3", defPerms}},
				})
			})
		})

		Convey("When search for hop3*", func() {
			found, _ := theHops.searchHop("hop3*")

			Convey("It should return only hop3", func() {
				So(found, ShouldResemble, &hops{
					"hop3": []Hop{&Docker{"a", "hop3", defPerms}},
				})
			})
		})

		Convey("When search for *", func() {
			found, _ := theHops.searchHop("*")

			Convey("It should return all hops", func() {
				So(found, ShouldResemble, &hops{
					"hop1": []Hop{&Docker{"a", "hop1", defPerms}},
					"hop3": []Hop{&Docker{"a", "hop3", defPerms}},
					"some": []Hop{&Docker{"a", "some", defPerms}},
				})
			})
		})

		Convey("When search for *o*", func() {
			found, _ := theHops.searchHop("*o*")

			Convey("It should return hop1, hop3, some", func() {
				So(found, ShouldResemble, &hops{
					"hop1": []Hop{&Docker{"a", "hop1", defPerms}},
					"hop3": []Hop{&Docker{"a", "hop3", defPerms}},
					"some": []Hop{&Docker{"a", "some", defPerms}},
				})
			})
		})

		Convey("When search for o*", func() {
			_, err := theHops.searchHop("o*")

			Convey("It should fail with NotFound error", func() {
				So(err, ShouldEqual, HopNotFound)
			})
		})
	})
}
