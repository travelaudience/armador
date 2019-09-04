package armador

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/travelaudience/armador/internal/commands"
)

func TestChartList_addFromArmadorFile(t *testing.T) {
	cases := []struct {
		name        string
		charts      ChartList
		armadorPath string
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// if err := tt.charts.addFromArmadorFile(tt.name, tt.armadorPath); (err != nil) != tt.wantErr {
			// 	t.Errorf("ChartList.addFromArmadorFile() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}

func TestChartList_mergeCharts(t *testing.T) {
	cases := []struct {
		name           string
		inputCharts    []Chart
		charts         ChartList
		expectedCharts ChartList
	}{
		{
			name: "merge-empty-chart",
			inputCharts: []Chart{
				Chart{},
			},
			charts: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedCharts: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
				"":            Chart{},
			},
		},
		{
			name: "merge-existing-chart",
			inputCharts: []Chart{
				Chart{Name: "secondchart", Repo: "stable"},
			},
			charts: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedCharts: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
		},
		{
			name: "merge-new-charts",
			inputCharts: []Chart{
				Chart{Name: "thirdchart", Repo: "stable"},
				Chart{Name: "fourth-chart", Repo: "stable"},
			},
			charts: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedCharts: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart":  Chart{Name: "secondchart", Repo: "stable"},
				"thirdchart":   Chart{Name: "thirdchart", Repo: "stable"},
				"fourth-chart": Chart{Name: "fourth-chart", Repo: "stable"},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.charts.mergeCharts(tt.inputCharts)
			if !reflect.DeepEqual(tt.charts, tt.expectedCharts) {
				t.Errorf("unexpected results for\n  TEST: %s \n  EXPECTED: %v\n  GOT: %v", tt.name, tt.expectedCharts, tt.charts)
			}
		})
	}
}
func TestChartList_flattenInitialChartsMap(t *testing.T) {
	cases := []struct {
		name     string
		input    ChartList
		expected map[string]Chart
	}{

		{
			name:     "empty",
			input:    make(ChartList),
			expected: make(map[string]Chart),
		},
		{
			name: "simple",
			input: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expected: map[string]Chart{
				"first-chart": Chart{},
				"secondchart": Chart{},
			},
		},
		{
			name: "simple dependecy",
			input: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expected: map[string]Chart{
				"first-chart": Chart{},
				"secondchart": Chart{},
				"first-dep":   Chart{Name: "first-dep"},
			},
		},
		{
			name: "dependecy with version",
			input: ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable",
					Dependencies: []Chart{Chart{Name: "first-dep", Repo: "stable", Version: "1.2.3"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expected: map[string]Chart{
				"first-chart": Chart{},
				"secondchart": Chart{},
				"first-dep":   Chart{Name: "first-dep", Repo: "stable", Version: "1.2.3"},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.flattenInitialChartsMap()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("unexpected results for %s. Got: %v Expected: %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestChartList_processCharts(t *testing.T) {
	dirs := commands.Dirs{
		Tmp:   commands.TmpDirs{Root: "tp/tmp/", Extracted: "tp/tmp/extracted", Hold: "tp/tmp/hold"},
		Cache: commands.CacheDirs{Root: "tp/cache/", Charts: "tp/cache/charts"},
	}
	cmd := commands.CmdMock{}
	tests := []struct {
		name           string
		charts         *ChartList
		expectedCharts *ChartList
		wantErr        bool
	}{
		{
			name:           "empty",
			charts:         &ChartList{},
			expectedCharts: &ChartList{},
			wantErr:        false,
		},
		{
			name: "no-dependecy",
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"},
			},
			wantErr: false,
		},
		{
			name: "simple-dependecy",
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "first-dep", Repo: "github", PathToChart: "helm/first-dep"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "first-dep", Repo: "github", PathToChart: "helm/first-dep"}},
				},
				"first-dep": Chart{Name: "first-dep", Repo: "github",
					ChartPath: "tp/cache/charts/first-dep/helm/first-dep", Packaged: false, PathToChart: "helm/first-dep",
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"},
			},
			wantErr: false,
		},
		{
			name: "complex-dependecy",
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "first-dep", Repo: "github", Version: "1.0.0", Packaged: false, PathToChart: "helm/first-dep"}}, Packaged: false, PathToChart: "helm/first-chart",
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "first-dep", Repo: "github", Version: "1.0.0", Packaged: false, PathToChart: "helm/first-dep"}}, Packaged: false, PathToChart: "helm/first-chart",
				},
				"first-dep": Chart{Name: "first-dep", Repo: "github", Version: "1.0.0",
					ChartPath: "tp/cache/charts/first-dep/helm/first-dep", Packaged: false, PathToChart: "helm/first-dep",
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"},
			},
			wantErr: false,
		},
		{
			name: "duplicate-dependecy",
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
						Dependencies: []Chart{Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"}},
					}},
				},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"}},
				},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData",
					Dependencies: []Chart{Chart{Name: "first-chart", Repo: "stable", ChartPath: "../testData",
						Dependencies: []Chart{Chart{Name: "secondchart", Repo: "stable", ChartPath: "../testData"}},
					}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			depList := tt.charts.flattenInitialChartsMap()
			filterDuplicates := make(map[string]struct{})
			err := tt.charts.processCharts(cmd, depList, dirs, filterDuplicates)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected test error for %s: %v", tt.name, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("expected error, no error received for %s", tt.name)
			}
			if diff := deep.Equal(tt.charts, tt.expectedCharts); diff != nil {
				t.Errorf("%s test failed.\nDiff: %s \nGot: %v \nExpected: %v", tt.name, diff, tt.charts, tt.expectedCharts)
			}

		})
	}
}

func TestChartList_processRawValues(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		charts         *ChartList
		expectedCharts *ChartList
		expectedUnused []string
	}{
		{
			name:  "no-values",
			input: []string{},
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedUnused: []string{},
		},
		{
			name:  "values-for-each-chart",
			input: []string{"first-chart.image=v1.2.3", "secondchart.v1=foo", "secondchart.v2=bar"},
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", SetValues: []string{"image=v1.2.3"}},
				"secondchart": Chart{Name: "secondchart", Repo: "stable", SetValues: []string{"v1=foo", "v2=bar"}},
			},
			expectedUnused: []string{},
		},
		{
			name:  "values-dont-match",
			input: []string{"first-chart.image=v1.2.3", "v1=foo", "second.v2=bar"},
			charts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable"},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedCharts: &ChartList{
				"first-chart": Chart{Name: "first-chart", Repo: "stable", SetValues: []string{"image=v1.2.3"}},
				"secondchart": Chart{Name: "secondchart", Repo: "stable"},
			},
			expectedUnused: []string{"v1=foo", "second.v2=bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUnused := tt.charts.processRawValues(tt.input)
			if !reflect.DeepEqual(gotUnused, tt.expectedUnused) {
				t.Errorf("%s unused values failed. \nGot: %v \nExpected: %v", tt.name, gotUnused, tt.expectedUnused)
				return
			}
			if !reflect.DeepEqual(tt.charts, tt.expectedCharts) {
				t.Errorf("%s updated charts failed. \nGot: %v \nExpected: %v", tt.name, tt.charts, tt.expectedCharts)
			}
		})
	}
}
