package lib

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"os"
	osexec "os/exec"
	"sort"
	"strings"
)

type Assert struct {
	Desc   string `yaml:"desc"`
	Ref    string `yaml:"ref"`
	Expect string `yaml:"expect"`
	Fix    string `yaml:"fix"`
}

func (a Assert) String() string {
	return fmt.Sprintf("Assert{desc:'%s',ref:'%s',expect:'%s',fix:'%s'}", a.Desc, a.Ref, a.Expect, a.Fix)
}

type YamlEnvGoal struct {
	Cmd    string   `yaml:"cmd"`
	Args   []string `yaml:"args,omitempty"`
	Assert *Assert  `yaml:"assert,omitempty"`
	Desc   string   `yaml:"desc"`
}

type YamlGoal struct {
	Envs   *map[string]YamlEnvGoal `yaml:"envs,omitempty"`
	Cmd    string                  `yaml:"cmd,omitempty"`
	Args   []string                `yaml:"args,omitempty"`
	Assert *Assert                 `yaml:"assert,omitempty"`
	Desc   string                  `yaml:"desc,omitempty"`
}

type Command struct {
	Name   string
	Cmd    string
	Args   []string
	Assert *Assert
	Env    string
	Desc   string
}

func (c Command) Cli() string {
	if len(c.Args) == 0 {
		return c.Cmd
	} else {
		return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
	}
}

func (c Command) String() string {
	return fmt.Sprintf("Command{name:'%s',env:'%v',Cli:'%s',assert:'%s'}", c.Name, c.Env, c.Cli(), c.Assert)
}

type Commands struct {
	Commands []Command
}

func (c *Commands) get(name string) (*Command, bool) {
	for _, command := range c.Commands {
		if command.Name == name {
			if command.Env != "" {
				Fatal("%s goal is referenced as an assertion but is environment dependant."+
					"It is not supported yet. Make it as a simple alias for now.", name)
			}
			return &command, true
		}
	}
	return nil, false
}

func (c *Commands) GetWithEnv(name string, env string) (*Command, bool) {
	for _, command := range c.Commands {
		if command.Name == name {
			if command.Env != "" && env != "" {
				return &command, true
			} else if command.Env == "" {
				return &command, true
			}
		}
	}
	return nil, false
}

func (c *Commands) runAssertion(assert Assert) {
	Info("⌛ Check precondition: %s", assert.Desc)
	// TODO: env or !env?
	ref, exists := c.get(assert.Ref)

	if exists {
		out := strings.TrimSpace(getOutput(ref))
		if out != assert.Expect {
			msg := fmt.Sprintf(
				"Precondition failed: %s\n\tOutput:   \"%s\"\n\tExpected: \"%s\"\n\tCLI: %s",
				ref.Name,
				out,
				assert.Expect,
				ref.Cli(),
			)
			if assert.Fix != "" {
				msg += fmt.Sprintf("\n\tFix: %s", assert.Fix)
			}
			Fatal(msg)
		} else {
			Info("✅ Precondition: " + assert.Desc)
		}
	} else {
		Fatal("Unknown assertion ref: %s", assert.Ref)
	}
}

func (c *Commands) Exec(name string, env string) {

	command, exists := c.GetWithEnv(name, env)
	if exists {
		msg := fmt.Sprintf("🔨 Exec %s", command.Name)
		if env != "" {
			msg += " on " + env
		}
		Info(msg)
		if command.Assert != nil {
			c.runAssertion(*command.Assert)
		}

		cmd := osexec.Command(command.Cmd, command.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()

		if err != nil {
			Info(fmt.Sprint(err))
			// TODO: code from child program
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	} else {
		Fatal("No such command in goal.yaml: %s", name)
	}

}

func (c *Commands) Render() {
	Info("Available goals:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"goal", "Environment", "CLI", "Description", "Assertions"})
	table.SetRowLine(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
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
		tablewriter.Colors{tablewriter.Normal},
	)
	for _, cmd := range c.Commands {
		assertion := ""
		if cmd.Assert != nil {
			ref, exists := c.get(cmd.Assert.Ref)
			if exists {
				assertion = fmt.Sprintf("[%s] %s", ref.Name, cmd.Assert.Desc)
			} else {

			}
		}
		table.Append([]string{cmd.Name, cmd.Env, cmd.Cli(), cmd.Desc, assertion})
	}
	table.Render()
}

func getOutput(command *Command) string {
	cmd := osexec.Command(command.Cmd, command.Args...)
	var output bytes.Buffer
	cmd.Stdout = &output

	// TODO: handle
	_ = cmd.Run()

	return output.String()
}

func normalizeArgs(args []string) []string {
	if args == nil {
		return []string{}
	} else {
		return args
	}
}

func validateAssert(name string, env string, assert *Assert) {
	if assert == nil {
		return
	}
	if assert.Ref == "" {
		if env == "" {
			Fatal("Malformed goals. %s.assert.ref could not be empty", name)
		} else {
			Fatal("Malformed goals. %s.%s.assert.ref could not be empty", name, env)
		}
	}
}

func parseEnvCommands(name string, envs map[string]YamlEnvGoal) []Command {
	var commands []Command
	for env, envCommand := range envs {
		args := normalizeArgs(envCommand.Args)
		if envCommand.Cmd == "" {
			Fatal("Malformed goals. %s.%s.cmd could not be empty", name, env)
		}
		validateAssert(name, env, envCommand.Assert)
		commands = append(commands, Command{
			Name:   name,
			Cmd:    envCommand.Cmd,
			Args:   args,
			Desc:   envCommand.Desc,
			Assert: envCommand.Assert,
			Env:    env,
		})
	}
	return sortCommands(commands)
}

// ParseCommands from byte input (YAML)
func ParseCommands(bytes []byte) (*Commands, error) {

	rawCommands := map[string]YamlGoal{}
	if err := yaml.Unmarshal(bytes, &rawCommands); err != nil {
		return nil, err
	}
	var res []Command
	for name, command := range rawCommands {
		if command.Envs != nil {
			res = append(res, parseEnvCommands(name, *command.Envs)...)
		} else {
			validateAssert(name, "", command.Assert)
			args := normalizeArgs(command.Args)
			res = append(res, Command{
				Name:   name,
				Cmd:    command.Cmd,
				Args:   args,
				Desc:   command.Desc,
				Assert: command.Assert,
			})
		}
	}

	return &Commands{Commands: sortCommands(res)}, nil
}

func sortCommands(commands []Command) (sorted []Command) {
	sorted = append(sorted, commands...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Name < sorted[j].Name {
			return true
		}
		if sorted[i].Name > sorted[j].Name {
			return false
		}
		return sorted[i].Env < sorted[j].Env
	})
	return
}
