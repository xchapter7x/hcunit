package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type RenderCommand struct {
	Writer   io.Writer
	Template string   `short:"t" long:"template" description:"path to yaml template you would like to render"`
	Values   []string `short:"c" long:"values" description:"path to values file(s) you would like to use for rendering"`
}

func (s *RenderCommand) Execute(args []string) error {
	s.setDefaults()
	valuesConfig, err := mergeValues(s.Values)
	if err != nil {
		return fmt.Errorf("failed merging values files %w ", err)
	}

	renderedOutput, err := validateAndRender(s.Template, valuesConfig)
	if err != nil {
		return fmt.Errorf("error while rendering: %w", err)
	}

	for filename, renderedFile := range renderedOutput {
		fmt.Fprintf(s.Writer, "---\n#%s\n%v\n\n", filepath.Base(filename), renderedFile)
	}

	return nil
}

func (s *RenderCommand) setDefaults() {
	if s.Writer == nil {
		s.Writer = os.Stdout
	}
}
