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

type Goal struct {
	Name   string
	Cmd    string
	Args   []string
	Assert []Assertion
	Env    string
	Desc   string
}

func (c Goal) Cli() string {
	if len(c.Args) == 0 {
		return c.Cmd
	} else {
		return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
	}
}

func (c Goal) String() string {
	return fmt.Sprintf("Goal{name:'%s',env:'%v',Cli:'%s',assert:'%s'}", c.Name, c.Env, c.Cli(), c.Assert)
}

type Goals struct {
	Commands []Goal
}

func (c *Goals) get(name string) (*Goal, bool) {
	for _, command := range c.Commands {
		if command.Name == name {
			if command.Env != "" {
				Fatal("❗ %s goal is referenced as an assertion but is environment dependant."+
					"It is not supported yet. Make it as a simple alias for now.", name)
			}
			return &command, true
		}
	}
	return nil, false
}

func (c *Goals) GetWithEnv(name string, env string) (*Goal, bool) {
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

func (c *Goals) Exec(name string, env string) {

	command, exists := c.GetWithEnv(name, env)
	if exists {
		msg := fmt.Sprintf("🔨 Exec %s", command.Name)
		if env != "" {
			msg += " on " + env
		}
		Info(msg)
		for _, assert := range command.Assert {
			Info("⌛ Check precondition: %s", assert.describe())
			if err := assert.check(*c); err != nil {
				Fatal(err.Error())
			}
			Info("✅ Precondition: %s" + assert.describe())
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
		Fatal("❗ No such command in goal.yaml: %s", name)
	}

}

func (c *Goals) Render() {
	Info("Available goals:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"goal", "Environment", "CLI", "Description", "Assertions"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetAutoWrapText(false)
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
		var assertions []string

		for _, assert := range cmd.Assert {
			assertions = append(assertions, assert.describe())
		}
		table.Append([]string{cmd.Name, cmd.Env, cmd.Cli(), cmd.Desc, strings.Join(assertions, "\n")})
	}
	table.Render()
}

func getOutput(name string, args ...string) string {
	cmd := osexec.Command(name, args...)
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

func mkAssertions(args []YamlAssert) (assertions []Assertion) {
	if args == nil {
		return assertions
	} else {
		for _, assertion := range args {
			if assertion.Ref != "" {
				assertions = append(assertions, RefAssertion{
					Desc:   assertion.Desc,
					Ref:    assertion.Ref,
					Expect: assertion.Expect,
					Fix:    assertion.Fix,
				})
			} else if assertion.TerraformWorkspace != "" {
				assertions = append(assertions, TerraformWorkspaceAssertion{
					Expect: assertion.TerraformWorkspace,
				})
			} else if assertion.KubectlContext != "" {
				assertions = append(assertions, KubectlContextAssertion{
					Expect: assertion.KubectlContext,
				})
			} else if assertion.GcloudProject != "" {
				assertions = append(assertions, GcloudProjectAssertion{
					Expect: assertion.GcloudProject,
				})
			}
		}
		return assertions
	}
}

func validateAssert(goal string, env string, idx int, assert YamlAssert) {
	var err string
	if assert.Ref == "" && assert.TerraformWorkspace == "" && assert.KubectlContext == "" && assert.GcloudProject == "" {
		err = fmt.Sprintf("one of [%s] must be specified for asserion", strings.Join(availableAssertions, ", "))
	}
	if assert.Ref != "" && assert.Expect == "" {
		err = "for 'ref' assertions specify expected output in 'expect'"
	}

	if err == "" {
		return
	} else {
		if env == "" {
			Fatal(fmt.Sprintf("❗ Malformed %s.assert.%d: %s", goal, idx, err))
		} else {
			Fatal(fmt.Sprintf("❗ Malformed %s.%s.assert.%d: %s", goal, env, idx, err))
		}
	}
}

func parseEnvCommands(goal string, envs map[string]YamlEnvGoal) []Goal {
	var commands []Goal
	for env, envCommand := range envs {
		args := normalizeArgs(envCommand.Args)
		if envCommand.Cmd == "" {
			Fatal("❗ Malformed goals. %s.%s.cmd could not be empty", goal, env)
		}
		for idx, assert := range envCommand.Assert {
			validateAssert(goal, env, idx, assert)
		}
		commands = append(commands, Goal{
			Name:   goal,
			Cmd:    envCommand.Cmd,
			Args:   args,
			Desc:   envCommand.Desc,
			Assert: mkAssertions(envCommand.Assert),
			Env:    env,
		})
	}
	return sortCommands(commands)
}

// ParseCommands from byte input (YAML)
func ParseCommands(bytes []byte) (*Goals, error) {

	rawCommands := map[string]YamlGoal{}
	if err := yaml.Unmarshal(bytes, &rawCommands); err != nil {
		return nil, err
	}
	var res []Goal
	for name, command := range rawCommands {
		if command.Envs != nil {
			res = append(res, parseEnvCommands(name, *command.Envs)...)
		} else {
			for idx, assert := range command.Assert {
				validateAssert(name, "", idx, assert)
			}
			args := normalizeArgs(command.Args)
			res = append(res, Goal{
				Name:   name,
				Cmd:    command.Cmd,
				Args:   args,
				Desc:   command.Desc,
				Assert: mkAssertions(command.Assert),
			})
		}
	}

	return &Goals{Commands: sortCommands(res)}, nil
}

func sortCommands(commands []Goal) (sorted []Goal) {
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
