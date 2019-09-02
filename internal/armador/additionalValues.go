package armador

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ghodss/yaml"

	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

type AdditionalValues []AdditionalValue

type AdditionalValue struct {
	Repo string
	Path []string
}

func getAdditionalValues() (retValues AdditionalValues, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to parse additional in config file")
			retValues = nil
		}
	}()

	additionalValues := viper.Get("additionalValues").([]interface{})
	for _, val := range additionalValues {
		val := val.(map[interface{}]interface{})

		newVal := AdditionalValue{
			Repo: val["repo"].(string),
		}

		paths := val["path"].([]interface{})
		for _, p := range paths {
			newVal.Path = append(newVal.Path, p.(string))
		}

		retValues = append(retValues, newVal)
	}
	return retValues, err
}

// the below functions are "Values" related - maybe this all goes under `helm` ????

func UnmarshalCache(vip *viper.Viper) (returnConf map[string]string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unable to parse cache file: %s", r)
			returnConf = map[string]string{}
		}
	}()
	returnConf = map[string]string{}

	charts := vip.AllSettings()
	for key, obj := range charts {
		ver, ok := obj.(string)
		if !ok {
			ver = ""
		}
		returnConf[key] = ver
	}
	return returnConf, err
}

func (valSettings AdditionalValues) GetGlobalOverrideString(cmd commands.Command, overridesDir string) error {
	for i, vals := range valSettings {
		clonePath := filepath.Join(overridesDir, strconv.Itoa(i))
		err := commands.GitGet(cmd, vals.Repo, clonePath)
		if err != nil {
			logger.GetLogger().Warnf("Unable to obtain override files: %s", err)
		}
		for _, v := range vals.Path {
			overrideFile := filepath.Join(clonePath, v)
			valuesMap, err := readValuesFile(overrideFile)
			if err != nil {
				return err
			}
			if valuesMap == nil {
				logger.GetLogger().Warnf("Override file not found in path: %s", overrideFile)
				continue
			}
			err = saveValuesToFile(valuesMap, overridesDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//for each chart/top level key in the values file, create a new file with the chart contents
func saveValuesToFile(valuesMap map[string]interface{}, overridesDir string) error {
	logger := logger.GetLogger()
	for key, val := range valuesMap {
		cont, err := yaml.Marshal(val)
		if err != nil {
			return fmt.Errorf("unable to marshal %s config to YAML: %v", key, err)
		}
		err = ioutil.WriteFile(filepath.Join(overridesDir, key+".yaml"), cont, 0644)
		if err != nil {
			return fmt.Errorf("unable to save %s config to YAML: %v", key, err)
		}
		logger.Debugf("Saved %s override config to path %s", key, overridesDir)
	}
	return nil
}

func readValuesFile(filepath string) (map[string]interface{}, error) {
	// check if file exists
	info, err := os.Stat(filepath)
	if err != nil || info.IsDir() {
		// no file to read
		return nil, nil
	}

	valuesMap := map[string]interface{}{}
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.GetLogger().Infof("Couldn't read the file %s: %s", filepath, err)
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &valuesMap); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %s", filepath, err)
	}
	return valuesMap, nil
}
