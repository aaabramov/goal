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
	"test": {
		Desc:   "Run go tests",
		Cmd:    "go",
		Args:   []string{"test", "-v", "./..."},
		Assert: nil,
	},
	"terraform-workspace": {
		Desc:   "Current terraform workspace",
		Cmd:    "terraform",
		Args:   []string{"workspace", "show"},
		Assert: nil,
	},
	"terraform-apply": {
		Desc: "See https://github.com/aaabramov/goal/tree/master/examples/terraform",
		Envs: &map[string]lib.YamlEnvGoal{
			"dev": {
				Desc: "Terraform apply on dev",
				Cmd:  "terraform",
				Args: []string{"apply", "-var-file", "vars/dev.tfvars"},
				Assert: []lib.YamlAssert{
					{
						TerraformWorkspace: "dev",
					}},
			},
			"stage": {
				Desc: "Terraform apply on stage",
				Cmd:  "terraform",
				Args: []string{"apply", "-var-file", "vars/stage.tfvars"},
				Assert: []lib.YamlAssert{
					{
						TerraformWorkspace: "stage",
					},
				},
			},
		},
	},
	"k8s-apply": {
		Desc: "See https://github.com/aaabramov/goal/tree/master/examples/kubectl",
		Envs: &map[string]lib.YamlEnvGoal{
			"dev": {
				Desc: "kubectl apply on dev",
				Cmd:  "kubectl",
				Args: []string{"apply", "-f", "deployment.yaml"},
				Assert: []lib.YamlAssert{
					{
						KubectlContext: "gke_project_region_dev",
					},
					{
						Approval: "yes",
					},
				},
			},
			"stage": {
				Desc: "kubectl apply on stage",
				Cmd:  "kubectl",
				Args: []string{"apply", "-f", "deployment.yaml"},
				Assert: []lib.YamlAssert{
					{
						KubectlContext: "gke_project_region_stage",
					},
					{
						Approval: "yes",
					},
				},
			},
		},
	},
	"helm-upgrade": {
		Desc: "See https://github.com/aaabramov/goal/tree/master/examples/helm",
		Envs: &map[string]lib.YamlEnvGoal{
			"dev": {
				Desc: "helm upgrade on dev",
				Cmd:  "helm",
				Args: []string{"upgrade", "release-name", "-f", "values.yaml", "-f", "values/dev.yaml", "."},
				Assert: []lib.YamlAssert{
					{
						KubectlContext: "gke_project_region_dev",
					}},
			},
			"stage": {
				Desc: "helm upgrade on stage",
				Cmd:  "helm",
				Args: []string{"upgrade", "release-name", "-f", "values.yaml", "-f", "values/stage.yaml", "."},
				Assert: []lib.YamlAssert{
					{
						KubectlContext: "gke_project_region_stage",
					}},
			},
		},
	},
	"gcloud-ssh": {
		Desc: "See https://github.com/aaabramov/goal/tree/master/examples/gcloud",
		Envs: &map[string]lib.YamlEnvGoal{
			"dev": {
				Desc: "SSH to dev",
				Cmd:  "gcloud",
				Args: []string{"compute", "ssh", "dev-vm", "--zone=us-central1-c"},
				Assert: []lib.YamlAssert{
					{
						GcloudProject: "dev-project",
					}},
			},
			"stage": {
				Desc: "SSH to stage",
				Cmd:  "gcloud",
				Args: []string{"compute", "ssh", "stage-vm", "--zone=us-central1-c"},
				Assert: []lib.YamlAssert{
					{
						GcloudProject: "stage-project",
					}},
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
