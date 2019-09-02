package armador

import (
	"os"
	"testing"

	"github.com/go-test/deep"
)

func Test_saveValuesToFile(t *testing.T) {
	cases := []struct {
		name          string
		values        map[string]interface{}
		destDir       string
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
			destDir:       "../testData/extracted/savedVals/",
			expectedFiles: []string{"oneChart.yaml", "twoChart.yaml"},
			wantErr:       false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Setup, make sure the dest directory exists
			err := os.MkdirAll(tt.destDir, os.ModePerm)
			if err != nil {
				t.Errorf("%s had a problem in setting up the test: %s", tt.name, err)
			}
			// test saveValuesToFile()
			err = saveValuesToFile(tt.values, tt.destDir)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
			}
			for _, expectedFile := range tt.expectedFiles {
				if _, err := os.Stat(tt.destDir + expectedFile); os.IsNotExist(err) {
					t.Errorf("Test %s failed, file %s not found in %s", tt.name, expectedFile, tt.destDir)
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
