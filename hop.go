package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/EverythingMe/gofigure"
	"github.com/EverythingMe/gofigure/yaml"
)

type hop struct {
	Container   string
	Entrypoint  string
	Permissions permissions
}

type permissions struct {
	Cwd bool
}

type hops map[string]hop

func main() {
	cmdName := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]
	hop, err := getHop(cmdName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hop)

	cmdArgs := os.Args[1:]
	dockerArgs := []string{"run"}

	dockerArgs = append(dockerArgs, []string{"--entrypoint=" + hop.Entrypoint}...)

	if hop.Permissions.Cwd {
		cwd, _ := os.Getwd()
		dockerArgs = append(dockerArgs, []string{"-v", cwd + ":/hopper", "-w=/hopper"}...)
	}

	dockerArgs = append(dockerArgs, []string{"-t", hop.Container}...)

	dockerArgs = append(dockerArgs, cmdArgs...)

	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	fmt.Println(dockerArgs)

	cmd := exec.Command(binary, dockerArgs...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
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
