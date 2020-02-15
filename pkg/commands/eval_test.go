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
			values    []string
			policy    string
			failsWith error
			skip      bool
			verbose   bool
		}{
			{
				name:      "invalid policy path given",
				template:  "testdata/templates/something.yml",
				values:    []string{"testdata/values.yml"},
				failsWith: commands.InvalidPolicyPath,
			},
			{
				name:      "passing policy on a single template",
				template:  "testdata/templates/something.yml",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/passing/passing.rego",
				failsWith: nil,
			},
			{
				name:      "duplicate test hash",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/duplicate_keynames.rego",
				failsWith: commands.DuplicatePolicyFailure,
			},
			{
				name:      "passing policy on a template directory",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/passing/passing.rego",
				failsWith: nil,
			},
			{
				name:      "failing policy on a single template",
				template:  "testdata/templates/something.yml",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/failing/failing.rego",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "failing policy on a template directory",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/failing/failing.rego",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "multifile failing policy on a template directory",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/failing",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "multifile passing policy on a template directory",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/passing",
				failsWith: nil,
			},
			{
				name:      "has a properly structured input object",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/parse_input.rego",
				failsWith: nil,
			},
			{
				name:      "values.yml available in input",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/values_in_input.rego",
				failsWith: nil,
			},
			{
				name:      "templates available in input",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/templates_in_input.rego",
				failsWith: nil,
			},
			{
				name:      "supports assert[_] rule query",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/alternate_keyword.rego",
				failsWith: nil,
			},
			{
				name:      "supports mulitple values files",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml", "testdata/added_values.yml"},
				policy:    "testdata/policy/individuals/multiple_values.rego",
				failsWith: nil,
			},
			{
				name:      "supports mulitple values files last file wins",
				template:  "testdata/templates",
				values:    []string{"testdata/added_values.yml", "testdata/values.yml"},
				policy:    "testdata/policy/individuals/multiple_values.rego",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "should error when no query match in rego",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/no_keyword.rego",
				failsWith: commands.UnmatchedQuery,
			},
			{
				name:      "no passing assertions",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/no_passing_valid.rego",
				failsWith: commands.PolicyFailure,
			},
			{
				name:      "verbosity on success should print trace information",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml", "testdata/added_values.yml"},
				policy:    "testdata/policy/individuals/multiple_values.rego",
				failsWith: nil,
				verbose:   true,
			},
			{
				name:      "verbosity on failure should print trace information",
				template:  "testdata/templates",
				values:    []string{"testdata/values.yml"},
				policy:    "testdata/policy/individuals/no_passing_valid.rego",
				failsWith: commands.PolicyFailure,
				verbose:   true,
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
					Verbose:  tt.verbose,
				}
				err := evalCmd.Execute([]string{})
				if err != nil && !errors.Is(err, tt.failsWith) {
					t.Errorf("expected error:\n%v\ngot:\n%v", tt.failsWith, err)
				}

				if !tt.verbose && stdOut.Len() > 0 {
					t.Errorf("when verbose is off this should always be empty, but it contains %v bytes", stdOut.Len())
				}

				if tt.verbose && stdOut.Len() == 0 {
					t.Errorf("we expected to print verbose trace information, but it is empty")
				}

				if err == nil && tt.failsWith != nil {
					t.Errorf("expected a failing policy %w but no failures found", tt.failsWith)
				}
			})
		}
	})
}
