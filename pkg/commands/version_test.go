package commands_test

import (
	"testing"

	"github.com/xchapter7x/hcunit/pkg/commands"
)

func TestVersionCommand(t *testing.T) {
	t.Run("should never return an error", func(t *testing.T) {
		version := new(commands.VersionCommand)
		err := version.Execute([]string{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
