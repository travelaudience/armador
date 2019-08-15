package armador

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

type ArmadorConfig struct {
	AdditionalValues AdditionalValues
}
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

func (valSettings AdditionalValues) GetGlobalOverrideString(cmd commands.Command, overridesDir string) {
	for i, vals := range valSettings {
		clonePath := filepath.Join(overridesDir, strconv.Itoa(i))
		err := commands.GitGet(cmd, vals.Repo, clonePath)
		if err != nil {
			logger.GetLogger().Warnf("Unable to obtain override files: %s", err)
		}
		for _, v := range vals.Path {
			parseOverrideFile(filepath.Join(clonePath, v), overridesDir)
		}
	}
}

func parseOverrideFile(overrideFile, overridesDir string) error {
	logger := logger.GetLogger()
	configName := strings.TrimSuffix(filepath.Base(overrideFile), filepath.Ext(overrideFile))
	configPath := strings.TrimSuffix(overrideFile, filepath.Base(overrideFile))
	logger.Debugf("Parsing override file %s in path %s", configName, configPath)
	vip, err := ReadFileToViper(configName, configPath) // TODO: here the configName is different, it's called `values`
	if err != nil {
		return err
	}
	if vip == nil {
		logger.Warnf("Override file %s not found in path %s", configName, configPath)
		return nil
	}
	// for each chart in vip, create a new file with the contents
	for key, val := range vip.AllSettings() {
		cont, err := yaml.Marshal(val)
		if err != nil {
			logger.Warnf("unable to marshal %s config to YAML: %v", key, err)
			continue
		}
		err = ioutil.WriteFile(filepath.Join(overridesDir, key+".yaml"), cont, 0644)
		if err != nil {
			logger.Warnf("unable to save %s config to YAML: %v", key, err)
			continue
		}
		logger.Debugf("Saved %s override config to path %s", key, overridesDir)
	}
	return nil
}

func ReadFileToViper(configName, configPath string) (*viper.Viper, error) {
	new := viper.New()
	new.SetConfigName(configName) // name of config file (without extension)
	new.AddConfigPath(configPath) // path to look for the config file in
	err := new.ReadInConfig()     // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.GetLogger().Debug(err) // expected file not found
			return nil, nil
		}
		// Other errors reading the config file should be addressed
		logger.GetLogger().Warnf("Problem with config file at %s", configPath)
		return nil, err
	}
	return new, nil
}
