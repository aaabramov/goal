package lib

import (
	"errors"
	"fmt"
	"strings"
)

type Assertion interface {
	describe() string
	check(commands Goals) error
}

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
				"Precondition failed: %s\n\tOutput:   \"%s\"\n\tExpected: \"%s\"\n\tCLI: %s",
				ref.Name,
				out,
				a.Expect,
				ref.Cli(),
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

// TerraformWorkspaceAssertion checks current Terraform workspace by executing `terraform workspace show`
// and compares its output with Expect
type TerraformWorkspaceAssertion struct {
	Expect string
}

func (a TerraformWorkspaceAssertion) describe() string {
	return fmt.Sprintf("Check if on %s workspace", a.Expect)
}

func (a TerraformWorkspaceAssertion) check(c Goals) error {
	out := strings.TrimSpace(getOutput("terraform", "workspace", "show"))
	if out == a.Expect {
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"❌ Precondition failed: %s\n"+
					"\tExpected terraform workspace to be:   \"%s\"\n"+
					"\tActual terraform workspace:           \"%s\"\n"+
					"\tFix:                                  \"terraform workspace select %s\"",
				a.describe(),
				out,
				a.Expect,
				a.Expect,
			),
		)
	}
}
