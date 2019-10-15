package commands_test

import (
	"bytes"
	"testing"

	"github.com/xchapter7x/hcunit/pkg/commands"
)

var controlYaml string = `---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
  labels:
    heritage: "Tiller"
    release: "hcunit-name"
    component: "hcunit-name-hcunitcomp"
spec:
  rules:
    - host: hcunit.com
      http:
        paths:
          - backend:
              servicePort: 8500

`

func TestRenderCommand(t *testing.T) {
	t.Run("should render the given template using the given values", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			template string
			values   string
		}{
			{"template filepath", "testdata/templates/something.yml", "testdata/values.yml"},
			{"template dir path", "testdata/templates", "testdata/values.yml"},
		} {
			t.Run(tt.name, func(t *testing.T) {
				stdOut := new(bytes.Buffer)
				renderer := &commands.RenderCommand{
					Writer:   stdOut,
					Template: tt.template,
					Values:   tt.values,
				}
				err := renderer.Execute([]string{})
				if err != nil {
					t.Errorf("should not have errored:\n%v", err)
				}

				if stdOut.String() == "---\n\n\n" {
					t.Errorf(
						"expected a rendered yaml got:\n'%s'",
						stdOut.String(),
					)
				}

				if stdOut.String() != controlYaml {
					t.Errorf(
						"rendered yaml is wrong. \nwanted:\n'%s'\n got:\n'%s'",
						controlYaml,
						stdOut.String(),
					)
				}
			})
		}
	})

	t.Run("should validate template & values paths", func(t *testing.T) {
		for _, tt := range []struct {
			name        string
			render      *commands.RenderCommand
			shouldError bool
		}{
			{
				name:        "no template or values",
				render:      &commands.RenderCommand{},
				shouldError: true,
			},
			{
				name:        "no template",
				render:      &commands.RenderCommand{Values: "hi/there"},
				shouldError: true,
			},
			{
				name:        "no values",
				render:      &commands.RenderCommand{Template: "hi/There"},
				shouldError: true,
			},
			{
				name:        "invalid template",
				render:      &commands.RenderCommand{Values: "hi/there", Template: "yo/yo"},
				shouldError: true,
			},
			{
				name:        "invliad values",
				render:      &commands.RenderCommand{Template: "hi/There", Values: "yo/yo"},
				shouldError: true,
			},
			{
				name:        "values is not a file",
				render:      &commands.RenderCommand{Template: "testdata/templates/something.yml", Values: "testdata/"},
				shouldError: true,
			},
			{
				name:        "valid template & values file paths",
				render:      &commands.RenderCommand{Template: "testdata/templates/something.yml", Values: "testdata/values.yml"},
				shouldError: false,
			},
			{
				name:        "valid template dir & values file path",
				render:      &commands.RenderCommand{Template: "testdata/templates", Values: "testdata/values.yml"},
				shouldError: false,
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.render.Execute([]string{})
				if err == nil && tt.shouldError {
					t.Errorf("we should have errored but we didnt")
				}

				if err != nil && !tt.shouldError {
					t.Errorf("error: %v", err)
				}
			})
		}
	})
}
