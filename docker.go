package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/fsouza/go-dockerclient"
)

type Docker struct {
	Image       string
	Command     string
	Permissions permissions
}

type permissions struct {
	Cwd bool
}

func (h *Docker) Run(cmdArgs ...string) (int, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	containerConfig := &docker.Config{
		Image:      h.Image,
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
