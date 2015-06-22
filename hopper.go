package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kardianos/osext"
	"github.com/mitchellh/go-homedir"
)

type Hop interface {
	Run(cmdArgs ...string) (int, error)
}

type hops map[string][]Hop

type workspace struct {
	Hops       hops
	BinDir     string
	HopperPath string
}

var (
	localHopsFile = "hop.yaml"
)

func buildWorkspace() (*workspace, error) {
	var hopsFile, binDir string
	var wsp workspace
	if inLocalMode() {
		hopsFile = localHopsFile
		binDir = "./"
	} else {
		var err error
		hopsFile, err = homedir.Expand("~/.hopper/hops/hop.yaml")
		if err != nil {
			log.Fatal(err)
		}
		binDir, err = homedir.Expand("~/.hopper/bin")
		if err != nil {
			log.Fatal(err)
		}
	}
	hopperPath, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}
	wsp = workspace{Hops: nil, BinDir: binDir, HopperPath: hopperPath}
	hs, err := LoadHops(hopsFile)
	if err != nil {
		return nil, err
	}
	wsp.Hops = hs
	return &wsp, nil
}

func inLocalMode() bool {
	_, err := os.Stat(localHopsFile)
	return err == nil
}

func (w *workspace) runHop(name string, args []string) {
	h, err := w.getHop(name)
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

func (w *workspace) getHop(name string) (Hop, error) {
	if h, exist := w.Hops[name]; exist {
		return h[0], nil
	} else {
		return nil, fmt.Errorf("Cannot find hop definition for: %q", name)
	}
}
