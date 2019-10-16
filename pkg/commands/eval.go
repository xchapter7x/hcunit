package commands

import (
	"fmt"
	"io"
)

type EvalCommand struct {
	Writer   io.Writer
	Template string `short:"t" long:"template" description:"path to yaml template you would like to render"`
	Values   string `short:"v" long:"values" description:"path to values file you would like to use for rendering"`
	Policy   string `short:"p" long:"policy" description:"path to rego policies to evaluate against rendered templates"`
}

func (s *EvalCommand) Execute(args []string) error {
	return fmt.Errorf("not yet implemented")
}
