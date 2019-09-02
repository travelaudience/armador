package helm

import (
	"testing"

	"github.com/travelaudience/armador/internal/commands"
)

func TestInstall(t *testing.T) {
	tests := []struct {
		name           string
		cmd            commands.Command
		chartPath      string
		namespace      string
		overridePath   string
		overrides      []string
		setValues      []string
		expectedResult string
		wantErr        bool
	}{
		{
			name:           "basic-install",
			cmd:            commands.CmdMock{},
			chartPath:      "../",
			namespace:      "test",
			expectedResult: "Dir: ../, Cmd: helm upgrade --install basic-install-test . --namespace test",
			wantErr:        false,
		},
		{
			name:           "with-overrides",
			cmd:            commands.CmdMock{},
			chartPath:      "example",
			namespace:      "test",
			overridePath:   "../testData/overrideFiles",
			overrides:      []string{"other-file"},
			setValues:      []string{"example=test"},
			expectedResult: "Dir: example, Cmd: helm upgrade --install with-overrides-test . --namespace test -f other-file -f ../testData/overrideFiles/with-overrides.yaml --set example=test",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Install(tt.cmd, tt.name, tt.chartPath, tt.namespace, tt.overridePath, tt.overrides, tt.setValues)
			if err != nil && !tt.wantErr {
				t.Errorf("%s: unexpected test error: %v", tt.name, err)
				return
			} else if err == nil && tt.wantErr {
				t.Errorf("%s: expected error, no error received", tt.name)
				return
			}

			if got != tt.expectedResult {
				t.Errorf("unexpected test failure:\n  TEST: %s \n  EXPECTED: %s\n  GOT: %s", tt.name, tt.expectedResult, got)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name           string
		cmd            commands.Command
		chartPath      string
		namespace      string
		overridePath   string
		overrides      []string
		setValues      []string
		expectedResult string
		wantErr        bool
	}{
		{
			name:           "basic-diff",
			cmd:            commands.CmdMock{},
			chartPath:      "../",
			namespace:      "test",
			expectedResult: "Dir: ../, Cmd: helm diff upgrade basic-diff-test . --allow-unreleased",
			wantErr:        false,
		},
		{
			name:           "with-overrides",
			cmd:            commands.CmdMock{},
			chartPath:      "example",
			namespace:      "test",
			overridePath:   "../testData",
			overrides:      []string{"other-file"},
			setValues:      []string{"example=test"},
			expectedResult: "Dir: example, Cmd: helm diff upgrade with-overrides-test . --allow-unreleased -f other-file -f ../testData/with-overrides.yaml --set example=test",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Diff(tt.cmd, tt.name, tt.chartPath, tt.namespace, tt.overridePath, tt.overrides, tt.setValues)
			if err != nil && !tt.wantErr {
				t.Errorf("%s: unexpected test error: %v", tt.name, err)
				return
			} else if err == nil && tt.wantErr {
				t.Errorf("%s: expected error, no error received", tt.name)
				return
			}

			if got != tt.expectedResult {
				t.Errorf("unexpected test failure:\n  TEST: %s \n  EXPECTED: %s\n  GOT:      %s", tt.name, tt.expectedResult, got)
			}
		})
	}
}

// func TestFetch(t *testing.T) {
// 	type args struct {
// 	}
// 	tests := []struct {
// 		name           string
// 		cmd            commands.Command
// 		chart          string
// 		repo           string
// 		version        string
// 		holdDir        string
// 		extractDir     string
// 		cacheDir       string
// 		expectedResult string
// 		wantErr        bool
// 	}{
// 		{
// 			name:           "",
// 			cmd:            commands.CmdMock{},
// 			chart:          "",
// 			repo:           "",
// 			version:        "",
// 			holdDir:        "../testData",
// 			extractDir:     "",
// 			cacheDir:       "",
// 			expectedResult: "",
// 			wantErr:        false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := Fetch(tt.cmd, tt.chart, tt.repo, tt.version, tt.holdDir, tt.extractDir, tt.cacheDir)
// 			if err != nil && !tt.wantErr {
// 				t.Errorf("%s: unexpected test error: %v", tt.name, err)
// 				return
// 			} else if err == nil && tt.wantErr {
// 				t.Errorf("%s: expected error, no error received", tt.name)
// 				return
// 			}

// 			if got != tt.expectedResult {
// 				t.Errorf("unexpected test failure:\n  TEST: %s \n  EXPECTED: %s\n  GOT: %s", tt.name, tt.expectedResult, got)
// 			}
// 		})
// 	}
// }
