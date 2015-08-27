// Main hopper logic and definitions
package main

import (
	"fmt"
	"os"

	"github.com/kardianos/osext"
	"github.com/mitchellh/go-homedir"
)

// Hop defines environment to run command from Hopper
type Hop interface {
	// Runs hop with given args
	Run(cmdArgs ...string) (int, error)
}

// Maps hop names to hop definitions
type hops map[string][]Hop

// Hopper workspace - depends from user or local mode
type workspace struct {
	Hops       hops
	BinDir     string
	HopperPath string
}

var (
	// hops definitions file
	localHopsFile = "hop.yaml"
)

// Sets Hopper workspace because it depends on
// run mode - user or project local.
func buildWorkspace() (*workspace, error) {
	var hopsFile, binDir string
	var wsp workspace
	if inLocalMode() {
		log.Debug("Hopper in local mode")
		hopsFile = localHopsFile
		binDir = "./"
	} else {
		log.Debug("Hopper in user mode")
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
	log.Debug("Hopper run from %v", hopperPath)
	wsp = workspace{Hops: nil, BinDir: binDir, HopperPath: hopperPath}
	hs, err := loadHops(hopsFile)
	if err != nil {
		return nil, err
	}
	wsp.Hops = hs
	return &wsp, nil
}

// Checks if we are in local mode - is there a local hops
// definition file?
func inLocalMode() bool {
	_, err := os.Stat(localHopsFile)
	return err == nil
}

// Runs named hop with given args in current workspace
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

// Gets hop for a given name, or fails if there's
// no such hop in current workspace
func (w *workspace) getHop(name string) (Hop, error) {
	if h, exist := w.Hops[name]; exist {
		return h[0], nil
	}
	return nil, fmt.Errorf("Cannot find hop definition for: %q", name)
}
