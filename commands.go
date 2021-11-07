package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	osexec "os/exec"
	"strings"
)

type Assert struct {
	Name   string `yaml:"name"`
	Ref    string `yaml:"ref"`
	Equals string `yaml:"equals"`
}

func (a Assert) String() string {
	return fmt.Sprintf("Assert{ref:'%s',equals:'%s'}", a.Ref, a.Equals)
}

type YamlCommand struct {
	Cmd    string   `yaml:"cmd"`
	Args   []string `yaml:"args"`
	Assert *Assert  `yaml:"assert"`
	Desc   string   `yaml:"desc"`
}

type Command struct {
	Name   string
	Cmd    string
	Args   []string
	Assert *Assert
	Desc   string
}

func (c Command) cli() string {
	return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
}

func (c Command) String() string {
	return fmt.Sprintf("Command{name:'%s',cli:'%s',assert:'%s'}", c.Name, c.cli(), c.Assert)
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

func (c *Commands) exec(name string) {

	command, exists := c.get(name)
	if exists {
		info("⚙️  Exec %s", command.Name)
		if command.Assert != nil {
			info("⌛ Check precondition: %s", command.Assert.Name)
			ref, exists := c.get(command.Assert.Ref)

			if exists {
				out := strings.TrimSpace(getOutput(ref))
				if out != command.Assert.Equals {
					fatal(
						"Precondition failed: %s\n\tOutput:   \"%s\"\n\tExpected: \"%s\"\n\tCLI: %s",
						ref.Name,
						out,
						command.Assert.Equals,
						ref.cli(),
					)
				} else {
					info("✅ Precondition: " + command.Assert.Name)
				}
			} else {
				fatal("Unknown assertion ref: %s", command.Assert.Ref)
			}
		}

		cmd := osexec.Command(command.Cmd, command.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()

		if err != nil {
			info(fmt.Sprint(err))
			// TODO: code from child program
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	} else {
		fatal("No such command in goal.yaml: %s", name)
	}

}

func (c *Commands) render() {
	info("Available commands:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "CLI", "Description", "Assertions"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.Normal},
		tablewriter.Colors{tablewriter.Normal},
		tablewriter.Colors{tablewriter.Normal},
	)
	for _, cmd := range c.commands {
		assertion := ""
		if cmd.Assert != nil {
			assertion = cmd.Assert.Name
		}
		table.Append([]string{cmd.Name, cmd.cli(), cmd.Desc, assertion})
		//fmt.Printf("\t%s: '%s' #%s\n", cmd.Name, cmd.cli(), cmd.Desc)
	}
	table.Render()
}
