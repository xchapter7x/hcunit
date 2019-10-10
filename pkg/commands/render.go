package commands

import "fmt"

type RenderCommand struct {
	Template string `short:"t" long:"template" description:"path to yaml template you would like to render"`
	Values   string `short:"v" long:"values" description:"path to values file you would like to use for rendering"`
}

func (x *RenderCommand) Execute(args []string) error {
	return fmt.Errorf("render is not yet supported")
}
