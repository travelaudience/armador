package armador

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/helm"
	"github.com/travelaudience/armador/internal/logger"
)

func Recreate(cmd commands.Cmd, cacheFile, namespace string, dirs commands.Dirs, dryRun bool) error {
	logger := logger.GetLogger()
	var wg sync.WaitGroup

	// read the cacheFile
	chartMap, err := readCacheFile(cacheFile)
	if err != nil {
		return err
	}
	if len(chartMap) < 1 {
		logger.Warn("The cache file was not parsed into any charts to install")
		return nil
	}

	// iterate over each chart
	wg.Add(len(chartMap))
	for chartName, version := range chartMap {
		go func(cmd commands.Cmd, chartName, version, namespace string, dirs commands.Dirs) {
			defer wg.Done()
			logger.Debugf("will install version %s of \t %s", version, chartName)
			// check if tar of that version is saved (extract)
			tarFilename, err := getTarFilename(filepath.Join(dirs.Cache.Charts, chartName, version))
			if err != nil || tarFilename == "" {
				logger.Warnf("version %s of %s is not downloaded", version, chartName)
				logger.Warn("---FOR NOW SKIP, CAUSE NEED TO DEAL WITH NON PACKAGED CHART VERSIONS---")
				return
			}
			// extract to tmp
			err = commands.Extract(filepath.Join(dirs.Cache.Charts, chartName, version), dirs.Tmp.Extracted, tarFilename)
			if err != nil {
				logger.Errorf("Extracting helm chart failed", err)
				return
			}
			// check override file
			overrideFile := filepath.Join(dirs.Cache.Overrides, chartName+".yaml")
			if _, err := os.Stat(overrideFile); os.IsNotExist(err) {
				logger.Warnf("Values file for %s does not exist", chartName)
				logger.Warn("---FOR NOW SKIP, PROBABLY JUST ERROR ---")
				return
			}
			chartPath := filepath.Join(dirs.Tmp.Extracted, chartName)
			if dryRun {
				logger.Debugf("Compare changes to: %s ", chartName)
				helm.DiffDirect(cmd, chartName, namespace, chartPath, dirs.Cache.Overrides)
			} else {
				logger.Debugf("Installing: %s ", chartName)
				helm.InstallDirect(cmd, chartName, namespace, chartPath, dirs.Cache.Overrides)
			}
		}(cmd, chartName, version, namespace, dirs)
	}
	wg.Wait()

	return nil
}

func readCacheFile(cacheFile string) (map[string]string, error) {
	vip := viper.New()
	vip.SetConfigFile(cacheFile)
	err := vip.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if vip == nil {
		return nil, nil
	}
	return UnmarshalCache(vip)
}

func getTarFilename(tarPath string) (string, error) {
	files, err := ioutil.ReadDir(tarPath)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tgz") {
			continue
		}
		return file.Name(), nil
	}
	return "", nil
}
