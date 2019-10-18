package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/helm/helm/pkg/renderutil"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var FilepathValueEmpty = errors.New("given filepath value is empty")
var FilepathDirUnexpected = errors.New("filepath given is a Dir. We expect a path to a file")
var InvalidPolicyPath = errors.New("invalid policy path")
var PolicyFailure = errors.New("your policy failed")
var expectQuery = regexp.MustCompile("^expect(_[a-zA-Z]+)*$")

func validateAndRender(template, values string) (map[string]string, error) {
	templateFiles, err := validateFileOrDirPath(template)
	if err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	valuesFile, err := validateFilePath(values)
	if err != nil {
		return nil, fmt.Errorf("values validation failed: %w", err)
	}

	return render(valuesFile, templateFiles)
}

func render(values io.ReadCloser, templates map[string]io.ReadCloser) (map[string]string, error) {
	var name string
	var reader io.ReadCloser
	var data []byte
	defer values.Close()
	chartTemplates := make([]*chart.Template, 0)
	for name, reader = range templates {
		defer reader.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		data = buf.Bytes()
		chartTemplates = append(chartTemplates, &chart.Template{Name: name, Data: data})
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(values)
	valuesRaw := buf.String()
	testChart := &chart.Chart{
		Metadata:  &chart.Metadata{Name: "hcunit"},
		Templates: chartTemplates,
		Values:    &chart.Config{Raw: valuesRaw},
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

func validateFileOrDirPath(filePath string) (map[string]io.ReadCloser, error) {
	if filePath == "" {
		return nil, FilepathValueEmpty
	}

	fileFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid Template path given: %w", err)
	}

	fileStatus, err := fileFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error while checking file status: %w", err)
	}

	fileMode := fileStatus.Mode()
	if fileMode.IsDir() {
		filePointers, err := fileFile.Readdir(-1)
		fileFile.Close()
		if err != nil {
			return nil, fmt.Errorf("reading files from directory failed: %w", err)
		}

		files := make(map[string]io.ReadCloser)

		for _, file := range filePointers {
			filePath := fmt.Sprintf("%s/%s", filePath, file.Name())
			fileReadCloser, err := os.Open(filePath)
			if err != nil {
				return nil, fmt.Errorf("reading file failed: %w", err)
			}

			files[filePath] = fileReadCloser
		}

		return files, nil
	}

	return map[string]io.ReadCloser{filePath: fileFile}, nil
}

func validateFilePath(filePath string) (*os.File, error) {
	if filePath == "" {
		return nil, FilepathValueEmpty
	}

	fileFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid Values path given: %w", err)
	}

	fileStatus, err := fileFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error while checking file status: %w", err)
	}

	fileMode := fileStatus.Mode()
	if fileMode.IsDir() {
		return nil, FilepathDirUnexpected
	}
	return fileFile, nil
}

func evalPolicyOnInput(writer io.Writer, policy string, namespace string, input interface{}) error {
	ctx := context.Background()
	buf := topdown.NewBufferTracer()
	_, err := rego.New(
		rego.Query(fmt.Sprintf("data.%s.%s", namespace, "expect[_]")),
		rego.Input(input),
		rego.Tracer(buf),
		rego.Load([]string{policy}, nil),
	).Eval(ctx)
	if err != nil {
		return fmt.Errorf("failed eval on policies: %w", err)
	}

	bufWriter := new(bytes.Buffer)
	topdown.PrettyTrace(bufWriter, *buf)
	fmt.Fprint(writer, bufWriter.String())
	if strings.Contains(bufWriter.String(), "Fail ") {
		return PolicyFailure
	}

	return nil
}
