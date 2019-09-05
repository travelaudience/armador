package armador

import (
	"testing"

	"github.com/go-test/deep"
)

func TestChart_parseArmadorFile(t *testing.T) {
	cases := []struct {
		name     string
		chart    Chart
		expected Chart
		wantErr  bool
	}{
		{
			name:     "no-file",
			chart:    Chart{Name: "first", Repo: "stable"},
			expected: Chart{Name: "first", Repo: "stable"},
			wantErr:  false,
		},
		{
			name:  "basic-file",
			chart: Chart{Name: "parsed-example", Repo: "stable", ChartPath: "../testData/basic", OverrideValueFiles: nil},
			expected: Chart{Name: "parsed-example", Repo: "stable", ChartPath: "../testData/basic", OverrideValueFiles: []string{"values-test.yaml"},
				Dependencies: []Chart{Chart{Name: "dep-chart", Repo: "test-stable", Version: "3.5.4", Packaged: true}},
			},
			wantErr: false,
		},
		{
			name:     "invalid-yaml-file",
			chart:    Chart{Name: "parsed-example", Repo: "stable", ChartPath: "../testData/invalid"},
			expected: Chart{Name: "parsed-example", Repo: "stable", ChartPath: "../testData/invalid"},
			wantErr:  true,
		},
		{
			name:     "non-packaged-chart",
			chart:    Chart{Name: "first", Repo: "stable", Packaged: false, PathToChart: "./chart"},
			expected: Chart{Name: "first", Repo: "stable", Packaged: false, PathToChart: "./chart"},
			wantErr:  false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chart.parseArmadorFile()
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
			}
			if diff := deep.Equal(tt.chart, tt.expected); diff != nil {
				t.Errorf("%s test failed.\nDiff: %s \nGot: \n%++v \nExpected: \n%++v", tt.name, diff, tt.chart, tt.expected)
			}
		})
	}
}
