package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/xchapter7x/hcunit/pkg/commands"
)

var Version = "0.0.0-localdev"
var Buildtime = "localdev-time"
var Platform = "localdev-platform"

type Options struct{}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func init() {
	parser.AddCommand(
		"render",
		"Render a template yaml",
		"Given a values and template it will render them to stdout",
		new(commands.RenderCommand),
	)
	parser.AddCommand(
		"version",
		"display version info",
		"display version information to stdout",
		&commands.VersionCommand{
			Version:   Version,
			Buildtime: Buildtime,
			Platform:  Platform,
		},
	)
	parser.AddCommand(
		"eval",
		"evaluate a policy on a chart + values",
		"given a OPA/Rego Policy one can evaluate if the rendered templates of a chart using a given values file meet the defined rules of the policy or not",
		new(commands.EvalCommand),
	)
}
