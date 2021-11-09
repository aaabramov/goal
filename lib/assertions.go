package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Assertion interface {
	describe() string
	check(commands Goals) error
}

var availableAssertions = []string{
	"ref",
	"terraform_workspace",
	"kubectl_context",
	"gcloud_project",
}

// === CUSTOM

// RefAssertion calls another goal and compares its trimmed output with Expect
type RefAssertion struct {
	Desc   string
	Ref    string
	Expect string
	Fix    string
}

func (a RefAssertion) describe() string {
	return a.Desc
}

func (a RefAssertion) check(c Goals) error {
	// TODO: env or !env?
	ref, exists := c.get(a.Ref)

	if exists {
		out := strings.TrimSpace(getOutput(ref.Name, ref.Args...))
		if out != a.Expect {
			msg := fmt.Sprintf(
				"Precondition failed: %s\n"+
					"\tOutput:   %s\n"+
					"\tExpected: %s\n"+
					"\tCLI:      %s",
				ref.Name,
				strconv.Quote(out),
				strconv.Quote(a.Expect),
				strconv.Quote(ref.Cli()),
			)
			if a.Fix != "" {
				msg += fmt.Sprintf("\n\tFix: %s", a.Fix)
			}
			return errors.New(msg)
		} else {
			Info("✅ Precondition: " + a.Desc)
			return nil
		}
	} else {
		return errors.New("Unknown assertion ref: " + a.Ref)
	}
}

// === TERRAFORM

// TerraformWorkspaceAssertion checks current Terraform workspace by executing `terraform workspace show`
// and compares its output with Expect
type TerraformWorkspaceAssertion struct {
	Expect string
}

func (a TerraformWorkspaceAssertion) describe() string {
	return fmt.Sprintf("terraform.workspace == %s", strconv.Quote(a.Expect))
}

func (a TerraformWorkspaceAssertion) check(_ Goals) error {
	out := strings.TrimSpace(getOutput("terraform", "workspace", "show"))
	if out == a.Expect {
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"❌ Precondition failed: %s\n"+
					"\tExpected terraform workspace to be: %s\n"+
					"\tActual terraform workspace:         %s\n"+
					"\tFix:                                \"terraform workspace select %s\"",
				a.describe(),
				strconv.Quote(out),
				strconv.Quote(a.Expect),
				a.Expect,
			),
		)
	}
}

// === KUBERNETES

// KubectlContextAssertion checks current `kubectl` context by executing `kubectl config current-context`
// and compares its output with Expect
type KubectlContextAssertion struct {
	Expect string
}

func (a KubectlContextAssertion) describe() string {
	return fmt.Sprintf("kubectl.context == %s", strconv.Quote(a.Expect))
}

func (a KubectlContextAssertion) check(_ Goals) error {
	out := strings.TrimSpace(getOutput("kubectl", "config", "current-context"))
	if out == a.Expect {
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"❌ Precondition failed: %s\n"+
					"\tExpected kubectl context to be: %s\n"+
					"\tActual kubectl context:         %s\n"+
					"\tFix:                            \"kubectl config use-context %s\"",
				a.describe(),
				strconv.Quote(out),
				strconv.Quote(a.Expect),
				a.Expect,
			),
		)
	}
}

// === GCP

// GcloudProjectAssertion checks current `gcloud` project by executing `gcloud config get-value project`
// and compares its output with Expect
type GcloudProjectAssertion struct {
	Expect string
}

func (a GcloudProjectAssertion) describe() string {
	return fmt.Sprintf("gcloud.project == %s", strconv.Quote(a.Expect))
}

func (a GcloudProjectAssertion) check(_ Goals) error {
	out := strings.TrimSpace(getOutput("gcloud", "config", "get-value", "project"))
	if out == a.Expect {
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"❌ Precondition failed: %s\n"+
					"\tExpected gcloud project to be: %s\n"+
					"\tActual gcloud project:         %s\n"+
					"\tFix:                           \"gcloud config set project %s\"",
				a.describe(),
				strconv.Quote(out),
				strconv.Quote(a.Expect),
				a.Expect,
			),
		)
	}
}
