package commands

import (
	"fmt"
	"io"
	"os"
)

type RenderCommand struct {
	Writer   io.Writer
	Template string `short:"t" long:"template" description:"path to yaml template you would like to render"`
	Values   string `short:"v" long:"values" description:"path to values file you would like to use for rendering"`
}

func (s *RenderCommand) Execute(args []string) error {
	s.setDefaults()
	renderedOutput, err := validateAndRender(s.Template, s.Values)
	if err != nil {
		return fmt.Errorf("error while rendering: %w", err)
	}

	for _, renderedFile := range renderedOutput {
		fmt.Fprintf(s.Writer, "---\n%v\n\n", renderedFile)
	}

	return nil
}

func (s *RenderCommand) setDefaults() {
	if s.Writer == nil {
		s.Writer = os.Stdout
	}
}
