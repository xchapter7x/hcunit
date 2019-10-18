package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type EvalCommand struct {
	Writer    io.Writer
	Template  string `short:"t" long:"template" description:"path to yaml template you would like to render"`
	Values    string `short:"c" long:"values" description:"path to values file you would like to use for rendering"`
	Policy    string `short:"p" long:"policy" description:"path to rego policies to evaluate against rendered templates"`
	Namespace string `short:"n" long:"namespace" description:"policy namespace to query for rules"`
	Verbose   bool   `short:"v" long:"verbose" description:"prints tracing output to stdout"`
}

func (s *EvalCommand) Execute(args []string) error {
	s.setDefaults()

	if s.Policy == "" {
		return InvalidPolicyPath
	}

	fileFile, err := os.Open(s.Policy)
	if err != nil {
		return InvalidPolicyPath
	}
	fileFile.Close()

	renderedOutput, err := validateAndRender(s.Template, s.Values)
	if err != nil {
		return fmt.Errorf("error while rendering: %w", err)
	}

	var valuesConfig interface{}
	valuesFile, err := ioutil.ReadFile(s.Values)
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(valuesFile, &valuesConfig)
	if err != nil {
		return fmt.Errorf("Unmarshal: %v", err)
	}

	policyInput := make(map[string]interface{})
	for fpath, template := range renderedOutput {
		var config interface{}
		err = yaml.Unmarshal([]byte(template), &config)
		if err != nil {
			return fmt.Errorf("Unmarshal: %v", err)
		}

		policyInput[fpath] = config
	}

	policyInput[s.Values] = valuesConfig
	return evalPolicyOnInput(s.Writer, s.Policy, s.Namespace, policyInput)
}

func (s *EvalCommand) setDefaults() {
	s.Writer = new(bytes.Buffer)
	if s.Verbose {
		s.Writer = os.Stdout
	}

	if s.Namespace == "" {
		s.Namespace = "main"
	}
}
