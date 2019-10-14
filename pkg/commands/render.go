package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/helm/helm/pkg/renderutil"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type RenderCommand struct {
	Writer   io.Writer
	Template string `short:"t" long:"template" description:"path to yaml template you would like to render"`
	Values   string `short:"v" long:"values" description:"path to values file you would like to use for rendering"`
}

func (s *RenderCommand) Execute(args []string) error {
	s.setDefaults()
	templateFile, err := validateTemplatePath(s.Template)
	if err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	valuesFile, err := validateValuesPath(s.Values)
	if err != nil {
		return fmt.Errorf("values validation failed: %w", err)
	}
	renderedOutput, err := render(
		valuesFile,
		map[string]io.Reader{
			s.Template: templateFile,
		},
	)
	if err != nil {
		return fmt.Errorf("error while rendering: %w", err)
	}

	for _, renderedFile := range renderedOutput {
		fmt.Fprintf(s.Writer, "---\n%v\n\n", renderedFile)
	}
	return nil
}

func render(values io.Reader, templates map[string]io.Reader) (map[string]string, error) {
	var name string
	var reader io.Reader
	var data []byte
	for name, reader = range templates {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		data = buf.Bytes()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(values)
	valuesRaw := buf.String()
	testChart := &chart.Chart{
		Metadata: &chart.Metadata{Name: "hcunit"},
		Templates: []*chart.Template{
			{Name: name, Data: data},
		},
		Values: &chart.Config{Raw: valuesRaw},
	}

	defaultConfig := &chart.Config{Raw: ""}
	defaultOptions := renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "hcunit-name",
			Time:      new(timestamp.Timestamp),
			Namespace: "hcunit-namespace",
			Revision:  1,
			IsUpgrade: false,
			IsInstall: true,
		},
	}
	return renderutil.Render(testChart, defaultConfig, defaultOptions)
}

func (s *RenderCommand) setDefaults() {
	if s.Writer == nil {
		s.Writer = os.Stdout
	}
}

func validateTemplatePath(templatePath string) (*os.File, error) {
	if templatePath == "" {
		return nil, fmt.Errorf("'Template' value is empty")
	}

	templateFile, err := os.Open(templatePath)
	if err != nil {
		return nil, fmt.Errorf("invalid Template path given: %w", err)
	}

	templateStatus, err := templateFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error while checking file status: %w", err)
	}

	templateMode := templateStatus.Mode()
	if templateMode.IsDir() {
		return nil, fmt.Errorf("Template directory not yet supported. Only a single file.")
	}

	return templateFile, nil
}

func validateValuesPath(valuesPath string) (*os.File, error) {
	if valuesPath == "" {
		return nil, fmt.Errorf("'Values' is empty")
	}

	valuesFile, err := os.Open(valuesPath)
	if err != nil {
		return nil, fmt.Errorf("invalid Values path given: %w", err)
	}

	valuesStatus, err := valuesFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error while checking file status: %w", err)
	}

	valuesMode := valuesStatus.Mode()
	if valuesMode.IsDir() {
		return nil, fmt.Errorf("Values path given is a Dir. We expect a path to a file")
	}
	return valuesFile, nil
}
