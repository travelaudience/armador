package armador

import (
	"os"
	"testing"
)

func Test_getChartFromFilename(t *testing.T) {
	cases := []struct {
		name            string
		armadorFilePath string
		wantName        string
		wantChartPath   string
		wantErr         bool
	}{
		{
			name:            "basic-test",
			armadorFilePath: "testData/basic/armador.yaml",
			wantName:        "basic",
			wantChartPath:   "/testData/basic/",
			wantErr:         false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			wd, _ := os.Getwd()
			chartPath := wd + tt.wantChartPath
			gotName, gotChartPath, err := getChartFromFilename(tt.armadorFilePath)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
				return
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("%s test failed. \ngotName = %v, \nwant %v", tt.name, gotName, tt.wantName)
			}
			if gotChartPath != chartPath {
				t.Errorf("%s test failed. \ngotChartPath = %v, \nwant %v", tt.name, gotChartPath, chartPath)
			}
		})
	}
}
