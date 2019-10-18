package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestCommands(t *testing.T) {
	gomega.RegisterTestingT(t)
	pathToCLI, err := gexec.Build("github.com/xchapter7x/hcunit/cmd/hcunit")
	defer gexec.CleanupBuildArtifacts()
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	t.Run("hcunit --help", func(t *testing.T) {
		command := exec.Command(pathToCLI, "--help")
		errOut := new(bytes.Buffer)
		stdOut := new(bytes.Buffer)
		session, err := gexec.Start(command, stdOut, errOut)
		if err != nil {
			t.Fatalf("failed running command: %v", err)
		}

		session.Wait(120 * time.Second)
		if session.ExitCode() > 0 {
			t.Errorf(
				"call failed: %v %v %v",
				session.ExitCode(),
				string(session.Out.Contents()),
				string(session.Err.Contents()),
			)
		}

		if !strings.HasPrefix(stdOut.String(), "Usage:") {
			t.Errorf(
				"expected help output. Instead got:\n%s",
				stdOut.String(),
			)
		}
	})

	t.Run("hcunit eval -t xxx -c xxx -p xxx -v", func(t *testing.T) {
		for _, tt := range []struct {
			name          string
			policy        string
			expectFailure bool
		}{
			{"failing policy should fail", "testdata/policy/failing", true},
			{"passing policy should pass", "testdata/policy/passing", false},
		} {
			t.Run(tt.name, func(t *testing.T) {
				command := exec.Command(
					pathToCLI,
					"eval",
					"-t", "testdata/templates/something.yml",
					"-c", "testdata/values.yml",
					"-p", tt.policy,
				)
				errOut := new(bytes.Buffer)
				stdOut := new(bytes.Buffer)
				session, err := gexec.Start(command, stdOut, errOut)
				if err != nil {
					t.Fatalf("failed running command: %v", err)
				}

				session.Wait(120 * time.Second)
				if tt.expectFailure && session.ExitCode() == 0 {
					t.Errorf(
						"this was expected to fail but did not: %v %v %v",
						session.ExitCode(),
						string(session.Out.Contents()),
						string(session.Err.Contents()),
					)
				}

				if !tt.expectFailure && session.ExitCode() > 0 {
					t.Errorf(
						"call failed unexpectedly: %v %v %v",
						session.ExitCode(),
						string(session.Out.Contents()),
						string(session.Err.Contents()),
					)
				}
			})
		}
	})

	t.Run("hcunit render -t xxx -v xxx", func(t *testing.T) {
		command := exec.Command(pathToCLI, "render", "-t", "testdata/templates/something.yml", "-c", "testdata/values.yml")
		errOut := new(bytes.Buffer)
		stdOut := new(bytes.Buffer)
		session, err := gexec.Start(command, stdOut, errOut)
		if err != nil {
			t.Fatalf("failed running command: %v", err)
		}

		session.Wait(120 * time.Second)
		if session.ExitCode() > 0 {
			t.Errorf(
				"call failed: %v %v %v",
				session.ExitCode(),
				string(session.Out.Contents()),
				string(session.Err.Contents()),
			)
		}

		controlYaml := `---
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
		if !strings.Contains(stdOut.String(), controlYaml) {
			t.Errorf(
				"expected output \n'%v'\n Instead got:\n%s",
				controlYaml,
				stdOut.String(),
			)
		}
	})

	t.Run("hcunit version", func(t *testing.T) {
		command := exec.Command(pathToCLI, "version")
		errOut := new(bytes.Buffer)
		stdOut := new(bytes.Buffer)
		session, err := gexec.Start(command, stdOut, errOut)
		if err != nil {
			t.Fatalf("failed running command: %v", err)
		}

		session.Wait(120 * time.Second)
		if session.ExitCode() != 0 {
			t.Errorf(
				"call failed: %v %v %v",
				session.ExitCode(),
				string(session.Out.Contents()),
				string(session.Err.Contents()),
			)
		}

		if !strings.Contains(stdOut.String(), Version) {
			t.Errorf(
				"expected version output. Instead got:\n%s",
				stdOut.String(),
			)
		}
		if !strings.Contains(stdOut.String(), Platform) {
			t.Errorf(
				"expected platform output. Instead got:\n%s",
				stdOut.String(),
			)
		}
		if !strings.Contains(stdOut.String(), Buildtime) {
			t.Errorf(
				"expected buildtime output. Instead got:\n%s",
				stdOut.String(),
			)
		}
	})
}
