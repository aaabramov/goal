package cmd

import (
	"github.com/aaabramov/goal/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create new goal.yaml file in current file",
	Long:  "Create new goal.yaml file in current file",
	Run: func(cmd *cobra.Command, args []string) {
		initGoals("goal.yaml")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	// TODO: support templates
}

var defaultGoals = map[string]lib.YamlGoal{
	"workspace": {
		Cmd:    "terraform",
		Args:   []string{"apply", "-var-file", "vars/dev.tfvars"},
		Assert: nil,
	},
	"tf-apply": {
		Envs: &map[string]lib.YamlEnvGoal{
			"dev": {
				Desc: "Terraform apply on dev",
				Cmd:  "terraform",
				Args: []string{"apply", "-var-file", "vars/dev.tfvars"},
				Assert: []lib.Assert{
					{
						Desc:   "Check if on dev workspace",
						Ref:    "workspace",
						Expect: "dev",
						Fix:    "terraform workspace select dev",
					}},
			},
			"stage": {
				Desc: "Terraform apply on stage",
				Cmd:  "terraform",
				Args: []string{"apply", "-var-file", "vars/stage.tfvars"},
				Assert: []lib.Assert{
					{
						Desc:   "Check if on stage workspace",
						Ref:    "workspace",
						Expect: "stage",
						Fix:    "terraform workspace select stage",
					},
				},
			},
		},
	},
}

func initGoals(filename string) {
	lib.Info("⌛ Generating default %s file", filename)
	bytes, err := yaml.Marshal(defaultGoals)
	if err != nil {
		lib.Fatal("❗ Failed to generate default YAML for goals")
	}
	if err = ioutil.WriteFile(filename, bytes, 0644); err != nil {
		lib.Fatal("❗ Failed to create %s", filename)
	}
	lib.Info("✅ Generated default %s file. Try running `goal` to see available goals.", filename)
}
