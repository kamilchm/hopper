package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "hopper"
)

func main() {
	initLogging()
	wrkspace, err := buildWorkspace()
	if err != nil {
		log.Fatal(err)
	}

	cmdName := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]

	if cmdName != appName {
		wrkspace.runHop(cmdName, os.Args[1:])
	}

	app := kingpin.New(appName, "Docker commands in your PATH")

	run := app.Command("run", "Run specified hop")
	runName := run.Arg("hop", "Hop name to run").Required().String()
	runArgs := run.Arg("args", "Arguments to hop").Required().Strings()

	install := app.Command("install",
		"Install hop.yaml as local commands")
	installForce := install.Flag("force",
		"Overwrite existing files").Short('f').Bool()
	installPattern := install.Arg("pattern",
		"Name pattern (* for any char)").Required().String()

	search := app.Command("search", "Search for hops")
	searchPattern := search.Arg("pattern",
		"Search pattern (* for any char)").Required().String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		wrkspace.runHop(*runName, *runArgs)
	case install.FullCommand():
		found, err := wrkspace.Hops.searchHop(*installPattern)
		if err != nil {
			log.Fatal(err)
		}

		for h := range *found {
			err := wrkspace.installHop(h, *installForce)
			if err != nil {
				log.Fatal(err)
			}
		}
	case search.FullCommand():
		found, err := wrkspace.Hops.searchHop(*searchPattern)
		if err != nil {
			log.Fatal(err)
		}

		for h := range *found {
			fmt.Println(h)
		}
	}
}
