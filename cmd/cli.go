package cmd

import (
	"fmt"
	"github.com/aaabramov/goal/lib"
	"strings"

	"github.com/spf13/cobra"
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli GOAL [--on env]",
	Short: "Show CLI for specific goal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		goal := args[0]
		if cmd, exists := commands.GetWithEnv(strings.TrimSpace(goal), env); exists {
			lib.Info(cmd.Cli())
		} else {
			msg := fmt.Sprintf("❗ No such goal: %s", goal)
			if env != "" {
				msg += fmt.Sprintf(" on env \"%s\"", env)
			}
			lib.Fatal(msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)

	cliCmd.Flags().StringVarP(&env, "on", "e", "", "Environment to use, example: goal cli tf-apply --on dev")
}
