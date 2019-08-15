package commands

import "testing"

func TestExtract(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		extractPath string
		fileName    string
		wantErr     bool
	}{
		{
			name:        "extract",
			filePath:    "../testData/",
			extractPath: "../testData/extracted/",
			fileName:    "testExample-0.1.0.tgz",
			wantErr:     false,
		},
		{
			name:        "false-file",
			filePath:    "../testData/",
			extractPath: "../testData/extracted/",
			fileName:    "non-exist.tgz",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Extract(tt.filePath, tt.extractPath, tt.fileName)
			if err != nil && !tt.wantErr {
				t.Errorf("%s: unexpected test error: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("%s: expected error, no error received", tt.name)
			}
		})
	}
}
