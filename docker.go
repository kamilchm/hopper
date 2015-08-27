// Docker interactions
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

// Docker Hop definition
type Docker struct {
	// Image where hop will be run
	Image string
	// Command to run hop in docker
	Command string
	// Hop permisions to local system
	Permissions permissions
}

// Defines hop permissions to local system
type permissions struct {
	Cwd bool
}

// Run hop with args as Docker container.
// Passes local stdin, print hop stdout and stderr.
func (d *Docker) Run(cmdArgs ...string) (int, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	containerConfig := &docker.Config{
		Image:      d.Image,
		Entrypoint: []string{d.Command},
		Cmd:        cmdArgs,
		Volumes:    make(map[string]struct{}),
		StdinOnce:  true,
	}

	// checks if there's something on stdin
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	// is it a piped stream?
	stdinPipe := false
	if fi.Mode()&os.ModeNamedPipe != 0 {
		stdinPipe = true
	}

	containerConfig.Tty = !stdinPipe
	containerConfig.OpenStdin = stdinPipe

	hostConfig := &docker.HostConfig{}

	// setup permissions to local system
	if d.Permissions.Cwd {
		hostWd, _ := os.Getwd()
		containerWd := "/hopper"
		containerConfig.WorkingDir = containerWd
		hostConfig.Binds = []string{hostWd + ":" + containerWd}
	}

	log.Debug("Creating container to run %v", d)
	container, err := client.CreateContainer(docker.CreateContainerOptions{"", containerConfig, hostConfig})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Debug("Removing temporary container %v", container)
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
	}()

	log.Debug("attaching stdin, stdout and stderr")
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

	log.Debug("starting container %v", container)
	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		log.Fatal(err)
	}

	return client.WaitContainer(container.ID)
}

const (
	stdinStream  byte = 0
	stdoutStream      = 1
	stderrStream      = 2
)

// SplitStream splits docker stream data into stdout and stderr.
// For specifications see http://goo.gl/Dnbcye
// TODO: check https://github.com/samalba/dockerclient/pull/3
func SplitStream(stream io.Reader, stdout, stderr io.Writer) error {
	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(stream, header); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("could not read header: %v", err)
		}

		var dest io.Writer
		switch header[0] {
		case stdinStream, stdoutStream:
			dest = stdout
		case stderrStream:
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
