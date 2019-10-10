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

	t.Run("hcunit render -t xxx -v xxx", func(t *testing.T) {
		command := exec.Command(pathToCLI, "render", "-t", "tpath", "-v", "vpath")
		errOut := new(bytes.Buffer)
		stdOut := new(bytes.Buffer)
		session, err := gexec.Start(command, stdOut, errOut)
		if err != nil {
			t.Fatalf("failed running command: %v", err)
		}

		session.Wait(120 * time.Second)
		if session.ExitCode() <= 0 {
			t.Errorf(
				"call failed: %v %v %v",
				session.ExitCode(),
				string(session.Out.Contents()),
				string(session.Err.Contents()),
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
