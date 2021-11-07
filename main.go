package main

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	osexec "os/exec"
	"strings"
)

func parseCommands(bytes []byte) (*Commands, error) {
	rawCommands := map[string]YamlCommand{}
	err := yaml.Unmarshal(bytes, &rawCommands)
	if err != nil {
		return nil, err
	}
	var res []Command
	for name, command := range rawCommands {
		var args []string
		if command.Args == nil {
			args = []string{}
		} else {
			args = command.Args
		}
		res = append(res, Command{
			Name:   name,
			Cmd:    command.Cmd,
			Args:   args,
			Desc:   command.Desc,
			Assert: command.Assert,
		})
	}
	return &Commands{commands: res}, nil
}

func main() {

	defaultFilename := "goal.yaml"

	if _, err := os.Stat(defaultFilename); errors.Is(err, os.ErrNotExist) {
		fatal("%s does not exist!", defaultFilename)
	}

	file, err := ioutil.ReadFile("goal.yaml")
	if err != nil {
		fatal("Failed to commands from %s file.", defaultFilename)
	}
	commands, _ := parseCommands(file)

	if len(os.Args) == 1 {
		commands.render()
	} else {
		// TODO: support several commands
		name := strings.TrimSpace(os.Args[1])
		commands.exec(name)
	}

}

func getOutput(command *Command) string {
	cmd := osexec.Command(command.Cmd, command.Args...)
	var output bytes.Buffer
	cmd.Stdout = &output

	// TODO: handle
	_ = cmd.Run()

	return output.String()
}
