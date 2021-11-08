/*
Copyright © 2021 Andrii Abramov

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

var goalFile string
var env string
var commands Commands

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goal [goal to run]",
	Short: "Define and safely run project scoped aliases",
	Long:  `Allows you to create local aliases withing directory/repository with proper assertions upon executions.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		res := make([]string, len(commands.commands))
		for _, command := range commands.commands {
			res = append(res, command.Name)
		}
		return res, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if env != "" {
			fmt.Println("Running on " + env)
		}
		if len(args) > 0 {
			commands.exec(strings.TrimSpace(args[0]), env)
		} else {
			commands.render()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(loadGoals)

	rootCmd.PersistentFlags().StringVarP(&goalFile, "config", "c", "goal.yaml", "goals file to use")
	rootCmd.Flags().StringVarP(&env, "on", "e", "", "Environment to use, example: goal tf-apply --on dev")
}

func normalizeArgs(args []string) []string {
	if args == nil {
		return []string{}
	} else {
		return args
	}
}

func parseEnvCommands(name string, envs map[string]YamlEnvGoal) (commands []Command) {
	for env, envCommand := range envs {
		args := normalizeArgs(envCommand.Args)
		commands = append(commands, Command{
			Name:   name,
			Cmd:    envCommand.Cmd,
			Args:   args,
			Desc:   envCommand.Desc,
			Assert: envCommand.Assert,
			Env:    env,
		})
	}
	return
}

// parseCommands from byte input (YAML)
func parseCommands(bytes []byte) error {

	rawCommands := map[string]YamlGoal{}
	if err := yaml.Unmarshal(bytes, &rawCommands); err != nil {
		return err
	}
	var res []Command
	for name, command := range rawCommands {
		if len(command.Envs) > 1 {
			res = append(res, parseEnvCommands(name, command.Envs)...)
		} else {
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
	commands = Commands{commands: res}
	return nil
}

func loadGoals() {
	if goalFile != "" {
		bytes, err := ioutil.ReadFile(goalFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read goal file:", goalFile)
		}
		if err := parseCommands(bytes); err != nil {
			fmt.Fprintln(os.Stderr, "Invalid goal file:", goalFile)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Empty goal file:", goalFile)
	}
}
