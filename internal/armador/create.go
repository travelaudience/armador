package armador

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/helm"
	"github.com/travelaudience/armador/internal/logger"
	"gopkg.in/yaml.v2"
)

const (
	// TODO: CONFIG: default repo should be set as global config
	repo = "stable"
	// TODO: CONFIG: file name should be part of a global config
	armadorFile = "armador.yaml"
)

func Create(cmd commands.Command, projectDir, namespace string, dirs commands.Dirs, dryRun bool, rawValues []string) error {
	logger := logger.GetLogger()
	var wg sync.WaitGroup
	charts := make(ChartList, 0)

	additionalValues, err := getAdditionalValues()
	if err != nil {
		logger.Warn(err)
	}
	preReqs, err := getPrereqCharts() // TODO: this needs to be swapped
	if err != nil {
		logger.Warn(err)
	}

	var armadorFiles []string
	var armadorFileErr, helmDiffErr error

	// update Helm repos
	wg.Add(1)
	go func(cmd commands.Command) {
		defer wg.Done()
		helm.RepoUpdate(cmd)
	}(cmd)

	if dryRun {
		// check that helm diff plugin is installed
		wg.Add(1)
		go func(cmd commands.Command) {
			defer wg.Done()
			helmDiffErr = helm.CheckPlugin(cmd, "diff")
		}(cmd)
	}

	//  get armador file(s) for project to be worked on
	wg.Add(1)
	go func(cmd commands.Command) {
		defer wg.Done()
		armadorFiles, armadorFileErr = cmd.Exec("find-armador-files", "find", projectDir, "-name", armadorFile)
	}(cmd)

	// set additional values
	wg.Add(1)
	go func(cmd commands.Command, valSettings AdditionalValues, dirs commands.Dirs) {
		defer wg.Done()
		valSettings.GetGlobalOverrideString(cmd, dirs.Tmp.Overrides)
	}(cmd, additionalValues, dirs)

	wg.Wait()

	if armadorFileErr != nil {
		return armadorFileErr
	}
	if helmDiffErr != nil {
		return helmDiffErr
	}

	if len(armadorFiles) < 1 {
		logger.Warnf(`There's no armador file available at: %s
    Without an Armador file, it's unclear what values/dependencies to use.
    Under most cases, if there is only one chart to install,
    you should just execute 'helm install'.

    If this is not the case than it's possibe the --projectDir (%s) may be wrong.`, projectDir, projectDir)
		return nil
	}

	for _, f := range armadorFiles {
		name, armadorPath, err := getChartFromFilename(strings.TrimSpace(f))
		if err != nil {
			return err
		}
		charts.addFromArmadorFile(name, armadorPath)
	}

	charts.mergeCharts(preReqs)

	logger.Debugf("Charts to analyze: \n\t %s", charts)

	filterDuplicates := make(map[string]struct{}, 0)
	err = charts.processCharts(cmd, charts.flattenInitialChartsMap(), dirs, filterDuplicates)
	if err != nil {
		return err
	}
	chartNames := make([]string, len(charts))
	i := 0
	for n := range charts {
		chartNames[i] = n
		i++
	}
	logger.Infof("%d charts to be installed: %v", len(chartNames), chartNames)

	if len(rawValues) > 0 {
		unusedValues := charts.processRawValues(rawValues)
		if len(unusedValues) > 0 {
			logger.Warnf("The following values passed were not used: %v", unusedValues)
		}
	}

	//  run helm installs (in parallel)
	wg.Add(len(charts))
	for _, name := range chartNames {
		go func(cmd commands.Command, name, namespace string, dirs commands.Dirs, chart Chart) {
			defer wg.Done()
			if dryRun {
				logger.Debugf("Compare changes to: %s ", name)
				output, err := helm.Diff(cmd, name, chart.ChartPath, namespace, dirs.Tmp.Overrides, chart.OverrideValueFiles, chart.SetValues)
				if err != nil {
					logger.Warnf("%s had problems: %s", name, err)
				} else if len(output) > 0 {
					logger.Infof("%s changes \n%s", name, output)
				} else {
					logger.Infof("No changes for %s", name)
				}
			} else {
				logger.Debugf("Installing: %s ", name)
				output, err := helm.Install(cmd, name, chart.ChartPath, namespace, dirs.Tmp.Overrides, chart.OverrideValueFiles, chart.SetValues)
				if err != nil {
					logger.Warnf("%s was not installed: %s", name, err)
				} else {
					logger.Infof("%s has been installed \n%s", name, output)
				}
			}
		}(cmd, name, namespace, dirs, charts[name])
	}
	wg.Wait()

	// create a cache file to be used by re-create
	if !dryRun {
		err = CreateSnapshot(cmd, namespace, dirs.Cache.Root, dirs.Cache.Overrides)
		if err != nil {
			return err
		}
	}

	return nil
}

func getChartFromFilename(armadorFilePath string) (name string, chartPath string, err error) {
	absPath, err := filepath.Abs(armadorFilePath)
	if err != nil {
		return "", "", err
	}
	logger.GetLogger().Debugf("The path to config file is: %s and \n\t the abs path is: %s", armadorFilePath, absPath)
	chartPath = strings.TrimSuffix(absPath, armadorFile)
	name = path.Base(chartPath)
	return name, chartPath, nil
}

// write values to cache file
func saveValuesToChartCacheFile(values, chartName, pathToSave string) error {
	appliedValuesPath := filepath.Join(pathToSave, chartName+".yaml")
	err := ioutil.WriteFile(appliedValuesPath, []byte(values), 0644)
	if err != nil {
		return err
	}
	return nil
}

// use `helm get values (chartname)` and save those to a file,
//   then save the structure of the chart name/version/path-to-those-values
//   and convert this structure to a single "cache" file
func CreateSnapshot(cmd commands.Command, namespace, pathToSave, valuesPath string) error {
	var wg sync.WaitGroup
	logger := logger.GetLogger()
	chartMap := map[string]string{}
	logger.Info("Creating a cache file for the charts installed....")

	// get list of all charts in namespace:
	charts, err := helm.ListCharts(cmd, namespace)
	if err != nil {
		return err
	}

	wg.Add(len(charts))
	for _, name := range charts {
		go func(cmd commands.Command, name string) {
			defer wg.Done()
			chartName := name
			vals, err := helm.CollectValues(cmd, chartName)
			if err != nil {
				logger.Warnf("Values for %s not obtained: %s", name, err)
				return
			}
			if strings.TrimSpace(vals) == "" || strings.TrimSpace(vals) == "{}" {
				logger.Infof("Values for %s don't need to be saved: %s", name, strings.TrimSpace(vals))
				return
			}
			ver, err := helm.GetChartVersion(cmd, chartName)
			if err != nil || ver == "" {
				logger.Warnf("Deployed version of %s not obtained: %s", name, err)
				return
			}
			err = saveValuesToChartCacheFile(vals, name, valuesPath)
			if err != nil {
				logger.Warnf("Cache file for %s not saved: %s", name, err)
				return
			}
			chartMap[name] = ver
		}(cmd, name)
	}
	wg.Wait()

	cont, err := yaml.Marshal(chartMap)
	if err != nil {
		return err
	}
	// TODO: CONFIG: make the name of cache file a global setting
	cacheFile := filepath.Join(pathToSave, "envCache.yaml")
	err = ioutil.WriteFile(cacheFile, cont, 0644)
	if err != nil {
		return err
	}
	logger.Infof("A cache for this release is available here: %s", cacheFile)
	return nil
}
