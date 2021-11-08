package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var env string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run GOAL [--on env]",
	Short: "Run specified goal",
	//Long:  `TODO`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		res := make([]string, len(commands.Commands))
		for _, command := range commands.Commands {
			res = append(res, command.Name)
		}
		return res, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			commands.Exec(strings.TrimSpace(args[0]), env)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&env, "on", "e", "", "Environment to use, example: goal tf-apply --on dev")
}
