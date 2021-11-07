package main

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	osexec "os/exec"
	"strings"
)

type YamlCommand struct {
	Cmd  string   `yaml:"cmd"`
	Args []string `yaml:"args"`
	Desc string   `yaml:"desc"`
}

type Command struct {
	Name string
	Cmd  string
	Args []string
	Desc string
}

func (c Command) cli() string {
	return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
}

type Commands struct {
	commands []Command
}

func (c *Commands) get(name string) (*Command, bool) {
	for _, command := range c.commands {
		if command.Name == name {
			return &command, true
		}
	}
	return nil, false
}

func (c *Commands) render() {
	fmt.Println("Available commands:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "CLI", "Description"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.Normal},
		tablewriter.Colors{tablewriter.Normal},
	)
	for _, cmd := range c.commands {
		table.Append([]string{cmd.Name, cmd.cli(), cmd.Desc})
		//fmt.Printf("\t%s: '%s' #%s\n", cmd.Name, cmd.cli(), cmd.Desc)
	}
	table.Render()
}

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
			Name: name,
			Cmd:  command.Cmd,
			Args: args,
			Desc: command.Desc,
		})
	}
	return &Commands{commands: res}, nil
}

func main() {

	defaultFilename := "goal.yaml"

	if _, err := os.Stat(defaultFilename); errors.Is(err, os.ErrNotExist) {
		errorFatal("%s does not exist!", defaultFilename)
	}

	file, err := ioutil.ReadFile("goal.yaml")
	if err != nil {
		errorFatal("Failed to commands from %s file.", defaultFilename)
	}
	commands, _ := parseCommands(file)

	if len(os.Args) == 1 {
		commands.render()
	} else {
		// TODO: support several commands
		name := strings.TrimSpace(os.Args[1])
		command, exists := commands.get(name)
		if exists {
			exec(command)
		} else {
			errorFatal("No such command in goal.yaml: %s", name)
		}
	}

}

func errorFatal(message string, args ...string) {
	msg := fmt.Sprintf(message, args)
	_, _ = os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}

func exec(command *Command) {
	cmd := osexec.Command(command.Cmd, command.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err))
		// TODO: code from child program
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
