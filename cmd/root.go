package cmd

import (
	"github.com/aaabramov/goal/lib"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var goalFile string
var commands *lib.Commands

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goal [goal to run]",
	Short: "Define and safely run project scoped aliases",
	Long:  `Allows you to create local aliases withing directory/repository with proper assertions upon executions.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		res := make([]string, len(commands.Commands))
		for _, command := range commands.Commands {
			res = append(res, command.Name)
		}
		return res, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			commands.Render()
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
}

func loadGoals() {
	if goalFile != "" {
		bytes, err := ioutil.ReadFile(goalFile)
		if err != nil {
			lib.Fatal("Failed to read goals file: %s\n"+
				"\t- check if goal.yaml files exists in current directory\n"+
				"\t- specify goal.yaml explicitly using -c flag, e.g. 'goal -c ../goal.yaml'\n"+
				"\t- run 'goal init' to generate example goal.yaml file in current directory", goalFile)
		}
		parsed, err := lib.ParseCommands(bytes)
		if err != nil {
			lib.Fatal("Invalid goals file: %s", goalFile)
		} else {
			commands = parsed
		}
	} else {
		lib.Fatal("Goals filename not specified. Either create goal.yaml file or specify location explicitly with -c option")
	}
}
