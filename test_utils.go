// Utils for use in tests
package main

import (
	"io/ioutil"
)

// Creates temporary dir than can be used for test
func createTestDir() string {
	emptyDir, err := ioutil.TempDir("", "hopper-test")
	if err != nil {
		log.Fatal(err)
	}
	return emptyDir
}
