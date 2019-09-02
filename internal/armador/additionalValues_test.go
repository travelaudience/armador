package armador

import (
	"testing"

	"github.com/go-test/deep"
)

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
