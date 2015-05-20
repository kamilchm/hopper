package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Hop interface {
	Run(cmdArgs ...string) (int, error)
}

type hops map[string][]Hop

func main() {
	cmdName := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]
	h, err := getHop(cmdName)
	if err != nil {
		log.Fatal(err)
	}

	cmdArgs := os.Args[1:]

	exitCode, err := h.Run(cmdArgs...)
	if err != nil {
		log.Fatal(err)
		os.Exit(exitCode)
	}

	os.Exit(exitCode)
}

func getHop(name string) (Hop, error) {
	localHops, err := LoadHops("hop.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if h, exist := localHops[name]; exist {
		return h[0], nil
	} else {
		return nil, fmt.Errorf("Cannot find hop definition for: %q", name)
	}
}
