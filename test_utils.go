package main

import (
	"io/ioutil"
)

func createTestDir() string {
	emptyDir, err := ioutil.TempDir("", "hopper-test")
	if err != nil {
		log.Fatal(err)
	}
	return emptyDir
}
