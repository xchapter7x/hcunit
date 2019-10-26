package commands_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/xchapter7x/hcunit/pkg/commands"
)

func TestEvalCommand(t *testing.T) {
	t.Run("given a successfully rendered template", func(t *testing.T) {
		for _, tt := range []struct {
			name      string
			template  string
			values    string
			policy    string
			failsWith error
			skip      bool
		}{
			{
				name:      "invalid policy path given",
				template:  "testdata/templates/something.yml",
				values:    "testdata/values.yml",
				failsWith: commands.InvalidPolicyPath,
			},
			{
				name:      "passing policy on a single template",
				template:  "testdata/templates/something.yml",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/passing/passing.rego",
				failsWith: nil,
			},
			{
				name:      "passing policy on a template directory",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/passing/passing.rego",
				failsWith: nil,
			},
			{
				name:      "failing policy on a single template",
				template:  "testdata/templates/something.yml",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/failing/failing.rego",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "failing policy on a template directory",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/failing/failing.rego",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "multifile failing policy on a template directory",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "multifile passing policy on a template directory",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/passing",
				failsWith: nil,
			},
			{
				name:      "has a properly structured input object",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/individuals/parse_input.rego",
				failsWith: nil,
			},
			{
				name:      "values.yml available in input",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/individuals/values_in_input.rego",
				failsWith: nil,
			},
			{
				name:      "templates available in input",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/individuals/templates_in_input.rego",
				failsWith: nil,
			},
			{
				name:      "supports assert[_] rule query",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/individuals/alternate_keyword.rego",
				failsWith: nil,
			},
			{
				name:      "should error when no query match in rego",
				template:  "testdata/templates",
				values:    "testdata/values.yml",
				policy:    "testdata/policy/individuals/no_keyword.rego",
				failsWith: commands.UnmatchedQuery,
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				if tt.skip {
					t.Skip(fmt.Sprintf("feature not implemented: %v", tt.name))
				}

				stdOut := new(bytes.Buffer)
				evalCmd := &commands.EvalCommand{
					Writer:   stdOut,
					Template: tt.template,
					Policy:   tt.policy,
					Values:   tt.values,
				}
				err := evalCmd.Execute([]string{})
				if err != nil && !errors.Is(err, tt.failsWith) {
					t.Errorf("expected error:\n%v\ngot:\n%v", tt.failsWith, err)
				}

				if err == nil && tt.failsWith != nil {
					t.Errorf("expected a failing policy %w but no failures found", tt.failsWith)
				}
			})
		}
	})
}
