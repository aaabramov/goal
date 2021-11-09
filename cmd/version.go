package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.11" // TODO: set during build

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the current version of the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("goal/" + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
