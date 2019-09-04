package armador

import (
	"strings"

	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/helm"
	"github.com/travelaudience/armador/internal/logger"
)

// ChartList the collection of charts that will be installed
type ChartList map[string]Chart

func (charts ChartList) addFromArmadorFile(name, armadorPath string) error {
	newChart := Chart{Name: name, ChartPath: armadorPath}
	err := newChart.parseArmadorFile()
	if err != nil {
		return err
	}
	charts[name] = newChart
	return nil
}

func (charts ChartList) mergeCharts(preReqs []Chart) {
	for _, c := range preReqs {
		name := c.Name
		if _, ok := charts[name]; ok {
			logger.GetLogger().Warnf("%s already exists. \nNew: \n%+v \nExisting: \n%+v", name, c, charts[name])
			continue
		}
		charts[name] = c
	}
}

func (charts ChartList) flattenInitialChartsMap() map[string]Chart {
	result := map[string]Chart{}
	for k, v := range charts {
		result[k] = Chart{}
		for _, d := range v.Dependencies {
			if _, ok := result[d.Name]; !ok {
				result[d.Name] = d
			}
		}
	}
	logger.GetLogger().Debugf("The starting structure of charts is: \n\t %v", result)
	return result
}

func (charts *ChartList) processCharts(cmd commands.Command, depList map[string]Chart, dirs commands.Dirs, filterDuplicates map[string]struct{}) error {
	for n, dep := range depList {
		// filterDuplicates is a map[string]struct{} because checking if string is in a slice is too much code....
		if _, ok := filterDuplicates[n]; ok {
			continue
		}
		chart := (*charts)[n]
		// sync chart values with dependency info
		if dep.Name != "" && chart.Name == "" {
			chart.Name = dep.Name
			chart.Version = dep.Version
			chart.Repo = repo // default helm repo
			if dep.Repo != "" {
				chart.Repo = dep.Repo
			}
			chart.Packaged = dep.Packaged
			chart.PathToChart = dep.PathToChart
		}

		if chart.ChartPath == "" {
			if chart.Packaged {
				chart.processHelm(cmd, dirs)
			} else {
				chart.processGit(cmd, dirs)
			}
		}

		_, err := helm.DepUpdate(cmd, chart.ChartPath)
		if err != nil {
			logger.GetLogger().Warnf("Unable to update dependecies for %s", chart.Name)
			return err
		}

		err = chart.parseArmadorFile()
		if err != nil {
			return err
		}
		for _, d := range chart.Dependencies {
			if _, ok := depList[d.Name]; !ok {
				depList[d.Name] = d
			}
		}
		delete(depList, n)
		filterDuplicates[n] = struct{}{}
		(*charts)[n] = chart
		err = charts.processCharts(cmd, depList, dirs, filterDuplicates)
		// if there was a failure, it's likely to persist all the other scans, so better to error out fast
		if err != nil {
			return err
		}
	}
	return nil
}

func (charts *ChartList) processRawValues(rawValues []string) []string {
	unusedValues := []string{}
	for _, v := range rawValues {
		chartName, valString := getChartFromValueString(v)
		if chartName == "" {
			unusedValues = append(unusedValues, v)
			continue
		}
		chart, ok := (*charts)[chartName]
		if !ok {
			unusedValues = append(unusedValues, v)
			continue
		}
		chart.SetValues = append(chart.SetValues, valString)
		(*charts)[chartName] = chart
	}
	return unusedValues
}

func getChartFromValueString(valStr string) (string, string) {
	pos := strings.Index(valStr, ".")
	if pos == -1 {
		return "", valStr
	}
	return valStr[0:pos], valStr[pos+1:]
}
