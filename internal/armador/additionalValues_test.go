package armador

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/go-test/deep"
)

func Test_saveValuesToFile(t *testing.T) {
	cases := []struct {
		name          string
		values        map[string]interface{}
		expectedFiles []string
		wantErr       bool
	}{
		{
			name: "two-charts",
			values: map[string]interface{}{
				"oneChart": map[string]interface{}{
					"values": "set",
				},
				"twoChart": map[string]interface{}{
					"OtherValues": true,
				},
			},
			expectedFiles: []string{"oneChart.yaml", "twoChart.yaml"},
			wantErr:       false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Setup, make sure the dest directory exists
			dir, err := ioutil.TempDir("", "armadorTest")
			if err != nil {
				t.Errorf("Setup for %s failed: %s", tt.name, err)
			}
			defer os.RemoveAll(dir) // clean up directory

			// test saveValuesToFile()
			err = saveValuesToFile(tt.values, dir)
			// check the error response is as expected
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
			}
			// for each value, check that a file exists, and the content is expected
			for key, val := range tt.values {
				newFile := key + ".yaml"
				savedValues := map[string]interface{}{}
				if _, err := os.Stat(filepath.Join(dir, newFile)); os.IsNotExist(err) {
					t.Errorf("Test %s failed, file %s not found in %s", tt.name, newFile, dir)
				}
				gotBytes, _ := ioutil.ReadFile(filepath.Join(dir, newFile))
				yaml.Unmarshal(gotBytes, &savedValues)
				if diff := deep.Equal(savedValues, val); diff != nil {
					t.Errorf("%s test failed:\nDiff: %s \nGot: %v \nExpected: %v", tt.name, diff, savedValues, val)
				}
			}

		})
	}
}

func Test_readValuesFile(t *testing.T) {
	cases := []struct {
		name     string
		filepath string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "empty-file",
			filepath: "../testData/overrideFiles/with-overrides.yaml",
			expected: map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "non-existant-file",
			filepath: "../testData/does-not-exist",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "folder",
			filepath: "../testData",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "simple-file",
			filepath: "../testData/overrideFiles/simple.yaml",
			expected: map[string]interface{}{
				"caseSensitive": "info",
				"is": map[string]interface{}{
					"valid": true,
				},
			},
			wantErr: false,
		},
		{
			name:     "invalid-yaml",
			filepath: "../testData/overrideFiles/invalid.yaml",
			expected: nil,
			wantErr:  true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readValuesFile(tt.filepath)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
			}
			if diff := deep.Equal(got, tt.expected); diff != nil {
				t.Errorf("%s test failed:\nDiff: %s \nGot: %v \nExpected: %v", tt.name, diff, got, tt.expected)
			}
		})
	}

}
