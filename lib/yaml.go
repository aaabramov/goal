package lib

import "fmt"

type YamlAssert struct {
	Desc               string `yaml:"desc"`
	Ref                string `yaml:"ref,omitempty"`
	Expect             string `yaml:"expect"`
	Fix                string `yaml:"fix"`
	TerraformWorkspace string `yaml:"terraform_workspace,omitempty"`
	KubectlContext     string `yaml:"kubectl_context,omitempty"`
	GcloudProject      string `yaml:"gcloud_project,omitempty"`
}

func (a YamlAssert) String() string {
	return fmt.Sprintf("YamlAssert{desc:'%s',ref:'%s',expect:'%s',fix:'%s'}", a.Desc, a.Ref, a.Expect, a.Fix)
}

type YamlEnvGoal struct {
	Cmd    string       `yaml:"cmd"`
	Args   []string     `yaml:"args,omitempty"`
	Assert []YamlAssert `yaml:"assert,omitempty"`
	Desc   string       `yaml:"desc"`
}

type YamlGoal struct {
	Envs   *map[string]YamlEnvGoal `yaml:"envs,omitempty"`
	Cmd    string                  `yaml:"cmd,omitempty"`
	Args   []string                `yaml:"args,omitempty"`
	Assert []YamlAssert            `yaml:"assert,omitempty"`
	Desc   string                  `yaml:"desc,omitempty"`
}
