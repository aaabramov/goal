package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	osexec "os/exec"
	"strings"
)

type Assert struct {
	Desc   string `yaml:"desc"`
	Ref    string `yaml:"ref"`
	Expect string `yaml:"expect"`
	Fix    string `yaml:"fix"`
}

func (a Assert) String() string {
	return fmt.Sprintf("Assert{ref:'%s',expect:'%s'}", a.Ref, a.Expect)
}

type YamlCommand struct {
	Cmd    string   `yaml:"cmd"`
	Args   []string `yaml:"args"`
	Assert *Assert  `yaml:"assert"`
	Env    string   `yaml:"env"`
	Desc   string   `yaml:"desc"`
}

type Command struct {
	Name   string
	Cmd    string
	Args   []string
	Assert *Assert
	Env    string
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
			info("⌛ Check precondition: %s", command.Assert.Desc)
			ref, exists := c.get(command.Assert.Ref)

			if exists {
				out := strings.TrimSpace(getOutput(ref))
				if out != command.Assert.Expect {
					msg := fmt.Sprintf(
						"Precondition failed: %s\n\tOutput:   \"%s\"\n\tExpected: \"%s\"\n\tCLI: %s",
						ref.Name,
						out,
						command.Assert.Expect,
						ref.cli(),
					)
					if command.Assert.Fix != "" {
						msg += fmt.Sprintf("\n\tFix: %s", command.Assert.Fix)
					}
					fatal(msg)
				} else {
					info("✅ Precondition: " + command.Assert.Desc)
				}
			} else {
				fatal("Unknown assertion ref: %s", command.Assert.Ref)
			}
		}

		cmd := osexec.Command(command.Cmd, command.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
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
	info("Available goals:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"goal", "CLI", "Description", "Assertions"})
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
			ref, exists := c.get(cmd.Assert.Ref)
			if exists {
				assertion = fmt.Sprintf("[%s] %s", ref.Name, cmd.Assert.Desc)
			} else {

			}
		}
		table.Append([]string{cmd.Name, cmd.cli(), cmd.Desc, assertion})
		//fmt.Printf("\t%s: '%s' #%s\n", cmd.Desc, cmd.cli(), cmd.Desc)
	}
	table.Render()
}
