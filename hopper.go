package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Hop interface {
	Run(cmdArgs ...string) (int, error)
}

type hops map[string][]Hop

const (
	appName  = "hopper"
	hopsFile = "hop.yaml"
)

func main() {
	cmdName := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]

	if cmdName != appName {
		runHop(cmdName, os.Args[1:])
	}

	app := kingpin.New(appName, "Docker commands in your PATH")

	run := app.Command("run", "Run specified hop")
	runName := run.Arg("hop", "Hop name to run").Required().String()
	runArgs := run.Arg("args", "Arguments to hop").Required().Strings()

	install := app.Command("install",
		"Install hop.yaml as local commands")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		runHop(*runName, *runArgs)
	case install.FullCommand():
		localHops, err := LoadHops(hopsFile)
		if err != nil {
			log.Fatal(err)
		}

		for h := range localHops {
			if err := os.Symlink(os.Args[0], h); err == nil {
				log.Println(h, "successfully installed")
			} else {
				log.Fatalf("Error while installing %v: %v", h, err)
			}
		}
	}
}

func runHop(name string, args []string) {
	h, err := getHop(name)
	if err != nil {
		log.Fatal(err)
	}

	exitCode, err := h.Run(args...)
	if err != nil {
		log.Fatal(err)
		os.Exit(exitCode)
	}

	os.Exit(exitCode)
}

func getHop(name string) (Hop, error) {
	localHops, err := LoadHops(hopsFile)
	if err != nil {
		log.Fatal(err)
	}

	if h, exist := localHops[name]; exist {
		return h[0], nil
	} else {
		return nil, fmt.Errorf("Cannot find hop definition for: %q", name)
	}
}
