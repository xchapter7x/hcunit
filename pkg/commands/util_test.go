package commands_test

import (
	"fmt"
	"testing"

	"github.com/xchapter7x/hcunit/pkg/commands"
)

func TestUnmarshalYamlMap(t *testing.T) {
	for _, tt := range []struct {
		name    string
		yamlMap map[string]string
		matcher func(map[string]interface{}) error
	}{
		{
			name:    "valid yaml should show up in unmarshalled output",
			yamlMap: map[string]string{"random.yml": "something: andvalue"},
			matcher: func(m map[string]interface{}) error {
				yamlObject := m["random.yml"].(map[string]interface{})
				unmarshalledValue := yamlObject["something"].(string)
				if unmarshalledValue != "andvalue" {
					return fmt.Errorf("unexpected values in unmarshalled object: %v", m)
				}

				return nil
			},
		},
		{
			name:    "non yaml files should be left as strings",
			yamlMap: map[string]string{"random.txt": "t// asdgon---9"},
			matcher: func(m map[string]interface{}) error {
				v := m["random.txt"]
				if v.(string) != "t// asdgon---9" {
					return fmt.Errorf("non yaml files should remain strings, instead: %#v", v)
				}

				return nil
			},
		},
		{
			name:    "empty yaml should not show up in unmarshalled output",
			yamlMap: map[string]string{"random.yml": ""},
			matcher: func(m map[string]interface{}) error {
				v, ok := m["random.yml"]
				if ok {
					return fmt.Errorf("empty file should not be in the unmarshalled map object: %#v", v)
				}

				return nil
			},
		},
		{
			name: "multidoc yaml should unmarshal into an array element for each doc",
			yamlMap: map[string]string{"random.yml": `---
something: firstdoc---

---

something: otherdoc`},
			matcher: func(m map[string]interface{}) error {
				v := m["random.yml"].([]interface{})
				if len(v) != 2 {
					return fmt.Errorf("multi-doc yaml was not unmarshalled properly: %v", v)
				}

				return nil
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inputObject, err := commands.UnmarshalYamlMap(tt.yamlMap)
			if err != nil {
				t.Errorf("unexpected error while unmarshalling: %w", err)
			}

			err = tt.matcher(inputObject)
			if err != nil {
				t.Errorf("unexpected error %w", err)
			}
		})
	}
}
