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
		initGoals()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	// TODO: support templates
}

func initGoals() {
	lib.Info("⌛ Generating default goal.yaml file")
	goals := map[string]lib.YamlGoal{
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
					Assert: &lib.Assert{
						Desc:   "Check if on dev workspace",
						Ref:    "workspace",
						Expect: "dev",
						Fix:    "terraform workspace select dev",
					},
				},
				"stage": {
					Desc: "Terraform apply on stage",
					Cmd:  "terraform",
					Args: []string{"apply", "-var-file", "vars/stage.tfvars"},
					Assert: &lib.Assert{
						Desc:   "Check if on stage workspace",
						Ref:    "workspace",
						Expect: "stage",
						Fix:    "terraform workspace select stage",
					},
				},
			},
		},
	}
	bytes, err := yaml.Marshal(goals)
	if err != nil {
		lib.Fatal("Failed to generate default YAML for goals")
	}
	if err = ioutil.WriteFile("goal.yaml", bytes, 0644); err != nil {
		lib.Fatal("Failed to create goal.yaml")
	}
	lib.Info("✅ Generated default goal.yaml file. Try running `goal` to see available goals.")
}
