package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/helm/helm/pkg/renderutil"
	"github.com/mitchellh/colorstring"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/tester"
	"github.com/open-policy-agent/opa/topdown"
	yaml "gopkg.in/yaml.v3"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var FilepathValueEmpty = errors.New("given filepath value is empty")
var FilepathDirUnexpected = errors.New("filepath given is a Dir. We expect a path to a file")
var UnmatchedQuery = errors.New("your given query did not yield any matches")
var InvalidPolicyPath = errors.New("invalid policy path")
var PolicyFailure = errors.New("your policy failed")
var DuplicatePolicyFailure = errors.New("duplicate rule names found")
var expectQuery = regexp.MustCompile("^expect(_[a-zA-Z]+)*$")

func mergeValues(valueFiles []string) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	for _, filePath := range valueFiles {
		currentMap := map[string]interface{}{}

		bytes, err := readFile(filePath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
		}
		base = mergeMaps(base, currentMap)
	}
	return base, nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func readFile(filePath string) ([]byte, error) {
	if strings.TrimSpace(filePath) == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filePath)
}

func validateAndRender(template string, values []string) (map[string]string, error) {
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

func UnmarshalYamlMap(in map[string]string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for fpath, template := range in {
		if filepath.Ext(fpath) == ".yml" || filepath.Ext(fpath) == ".yaml" {
			documents := strings.Split(template, "\n---\n")
			var configDocs []interface{}
			for _, doc := range documents {
				var config interface{}
				err := yaml.Unmarshal([]byte(doc), &config)
				if err != nil {
					return nil, fmt.Errorf("Unmarshal '%s' failed: %v", fpath, err)
				}

				if config != nil {
					configDocs = append(configDocs, config)
				}
			}

			if configDocs != nil && len(configDocs) > 1 {
				out[filepath.Base(fpath)] = configDocs
			}

			if configDocs != nil && len(configDocs) == 1 {
				out[filepath.Base(fpath)] = configDocs[0]
			}

		} else {
			out[filepath.Base(fpath)] = template
		}
	}
	return out, nil
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

func validateFilePath(filePaths []string) (io.ReadCloser, error) {
	if len(filePaths) == 0 {
		return nil, FilepathValueEmpty
	}

	for _, filePath := range filePaths {
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
	return nil, FilepathValueEmpty
}

func getQueryList(policy string) map[string]int {
	res := map[string]int{}
	mods, _, _ := tester.Load([]string{policy}, nil)
	for _, mod := range mods {
		for _, rule := range mod.Rules {
			if strings.HasPrefix("expect[", string(rule.Head.Name)) ||
				strings.HasPrefix("assert[", string(rule.Head.Name)) {
				res[fmt.Sprintf("%s[%s]", rule.Head.Name, rule.Head.Key)] += 1
			}
		}
	}
	return res
}

func evalPolicyOnInput(writer io.Writer, policy string, namespace string, input interface{}) error {
	bufWriter := new(bytes.Buffer)
	testResults := make(map[string]bool)
	ctx := context.Background()
	var results rego.ResultSet
	queryList := getQueryList(policy)
	for querySuffix, querymatches := range queryList {
		if querymatches > 1 {
			colorstring.Println("[red]ERROR: you are using duplicate test names or variables. This could cause test failures to NOT be detected properly")
			colorstring.Println(fmt.Sprintf("[yellow]DUPLICATE KEY: %s", querySuffix))
			return DuplicatePolicyFailure
		}

		queryString := fmt.Sprintf("data.%s.%s", namespace, querySuffix)
		buf := topdown.NewBufferTracer()
		r := rego.New(
			rego.Query(queryString),
			rego.Tracer(buf),
			rego.Load([]string{policy}, nil),
		)
		query, err := r.PrepareForEval(ctx)
		if err != nil {
			return fmt.Errorf("failed preparing for eval on policies: %w", err)
		}

		resultSet, err := query.Eval(ctx, rego.EvalInput(input))
		if err != nil {
			return fmt.Errorf("failed eval on policies: %w", err)
		}

		testResults[queryString] = false
		for _, result := range resultSet {

			for _, expression := range result.Expressions {
				if expression.Text == queryString {
					testResults[queryString] = true
				}
			}
		}

		if len(resultSet) > 0 {
			results = append(results, resultSet...)
			topdown.PrettyTrace(bufWriter, *buf)
			fmt.Fprint(writer, bufWriter.String())
		}
	}

	if len(results) <= 0 {
		return UnmatchedQuery
	}

	testFailed := false
	for testname, passed := range testResults {
		if passed {
			colorstring.Print("[green]PASS: ")
			fmt.Println(testname)
		} else {
			testFailed = true
			colorstring.Print("[red]FAIL: ")
			fmt.Println(testname)
		}
	}

	if testFailed {
		colorstring.Println("[_red_][FAILURE] Policy violations found on the Helm Chart!")
		return PolicyFailure
	}

	colorstring.Println("[green][SUCCESS] Your Helm Chart complies with all policies!")
	return nil
}
