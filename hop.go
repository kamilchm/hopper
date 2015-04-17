package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EverythingMe/gofigure"
	"github.com/EverythingMe/gofigure/yaml"

	"github.com/fsouza/go-dockerclient"
)

type hop struct {
	Docker      string
	Command     string
	Permissions permissions
}

func (h *hop) run(cmdArgs ...string) (int, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	containerConfig := &docker.Config{
		Image:      h.Docker,
		Entrypoint: []string{h.Command},
		Cmd:        cmdArgs,
		Volumes:    make(map[string]struct{}),
		StdinOnce:  true,
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		containerConfig.Tty = true
		containerConfig.OpenStdin = false
	} else {
		containerConfig.Tty = false
		containerConfig.OpenStdin = true
	}

	hostConfig := &docker.HostConfig{}

	if h.Permissions.Cwd {
		hostWd, _ := os.Getwd()
		containerWd := "/hopper"
		containerConfig.WorkingDir = containerWd
		hostConfig.Binds = []string{hostWd + ":" + containerWd}
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{"", containerConfig, hostConfig})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
	}()

	attachChan := make(chan struct{})
	go func(succChan chan struct{}) {
		outWr := bufio.NewWriter(os.Stdout)
		errWr := bufio.NewWriter(os.Stderr)
		defer outWr.Flush()
		defer errWr.Flush()
		err := client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    container.ID,
			Stdout:       true,
			Stderr:       true,
			Stdin:        true,
			OutputStream: outWr,
			ErrorStream:  errWr,
			InputStream:  os.Stdin,
			Stream:       true,
			RawTerminal:  true,
			Success:      succChan,
		})
		if err != nil {
			log.Fatal(err)
		}
	}(attachChan)

	_, ok := <-attachChan
	if ok {
		attachChan <- struct{}{}
	}

	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		log.Fatal(err)
	}

	return client.WaitContainer(container.ID)
}

type permissions struct {
	Cwd bool
}

type hops map[string]hop

func main() {
	cmdName := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]
	h, err := getHop(cmdName)
	if err != nil {
		log.Fatal(err)
	}

	cmdArgs := os.Args[1:]

	exitCode, err := h.run(cmdArgs...)
	if err != nil {
		log.Fatal(err)
		os.Exit(exitCode)
	}

	os.Exit(exitCode)
}

func getHop(name string) (hop, error) {
	var localHops = make(hops)

	loader := gofigure.NewLoader(yaml.Decoder{}, true)

	err := loader.LoadFile(&localHops, "hop.yaml")
	if err != nil {
		panic(err)
	}

	if h, exist := localHops[name]; exist {
		return h, nil
	} else {
		return h, fmt.Errorf("Cannot find hop definition for: %q", name)
	}
}
