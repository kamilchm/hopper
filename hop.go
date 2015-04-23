package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

	stdinPipe := false
	if fi.Mode()&os.ModeNamedPipe != 0 {
		stdinPipe = true
	}

	containerConfig.Tty = !stdinPipe
	containerConfig.OpenStdin = stdinPipe

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
		var outBuf bytes.Buffer
		var outWr, errWr *bufio.Writer
		if stdinPipe {
			outWr = bufio.NewWriter(&outBuf)
			errWr = bufio.NewWriter(&outBuf)
		} else {
			outWr = bufio.NewWriter(os.Stdout)
			errWr = bufio.NewWriter(os.Stderr)
		}
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

		if stdinPipe {
			SplitStream(bufio.NewReader(&outBuf), os.Stdout, os.Stderr)
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

const (
	StdinStream  byte = 0
	StdoutStream      = 1
	StderrStream      = 2
)

// SplitStream splits docker stream data into stdout and stderr.
// For specifications see http://goo.gl/Dnbcye
func SplitStream(stream io.Reader, stdout, stderr io.Writer) error {
	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(stream, header); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return fmt.Errorf("could not read header: %v", err)
			}
		}

		var dest io.Writer
		switch header[0] {
		case StdinStream, StdoutStream:
			dest = stdout
		case StderrStream:
			dest = stderr
		default:
			return fmt.Errorf("bad STREAM_TYPE given: %x", header[0])
		}

		frameSize := int64(binary.BigEndian.Uint32(header[4:]))
		if _, err := io.CopyN(dest, stream, frameSize); err != nil {
			return fmt.Errorf("copying stream payload failed: %v", err)
		}
	}
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
