package armador

import (
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/helm"
	"github.com/travelaudience/armador/internal/logger"
)

// Chart the base structure of all things related to a helm chart
type Chart struct {
	Name               string
	Repo               string // helm chart repo - if packaged=false, than this would be the git repo
	Version            string
	Dependencies       []Chart
	OverrideValueFiles []string `yaml:"overrideValueFiles"`
	SetValues          []string // these come from using `--set` in the cli
	ChartPath          string   // the local path to the extracted chart
	Packaged           bool     // does the chart come as a tar
	PathToChart        string   // if packaged=false, where in the repo is the chart located
}

func (chart *Chart) parseArmadorFile() error {
	// configName := strings.TrimSuffix(armadorFile, ".yaml")
	vip, err := readFileToViper("armador", chart.ChartPath) // TODO: CONFIG: the config name should be set somehow
	if err != nil {
		return err
	}
	if vip == nil {
		// there was no file to be parsed
		return nil
	}
	chart.unmarshalArmadorConfig(vip)
	logger.GetLogger().Debugf("%s chart at %s has the following structure: %+v", chart.Name, chart.ChartPath, chart)
	return nil
}

func (chart *Chart) processHelm(cmd commands.Command, dirs commands.Dirs) error {
	p, err := helm.Fetch(cmd, chart.Name, chart.Repo, chart.Version, dirs.Tmp.Hold, dirs.Tmp.Extracted, dirs.Cache.Charts)
	if err != nil {
		logger.GetLogger().Warnf("Unable to fetch and extract %s: %s", chart.Name, err)
		return err
	}
	chart.ChartPath = p
	return nil
}

// manage non-packaged charts
func (chart *Chart) processGit(cmd commands.Command, dirs commands.Dirs) error {
	clonePath := filepath.Join(dirs.Cache.Charts, chart.Name)
	err := commands.GitGet(cmd, chart.Repo, clonePath) // TODO: be able to specificy version to clone
	if err != nil {
		return err
	}
	chart.ChartPath = filepath.Join(clonePath, chart.PathToChart)
	return nil
}

func readFileToViper(configName, configPath string) (*viper.Viper, error) {
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

func (chart *Chart) unmarshalArmadorConfig(vip *viper.Viper) {
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Warnf("Unable to parse helm config info: %s", r)
		}
	}()

	ovf := []string{}
	ovfs, _ := vip.Get("overrideValueFiles").([]interface{})
	for _, f := range ovfs {
		ovf = append(ovf, f.(string))
	}
	chart.OverrideValueFiles = ovf

	// TODO: add ability for non-pacakaged charts to be dependencies (ie: use same logic as in `getPrereqCharts()`)
	deps, ok := vip.Get("dependencies").([]interface{})
	if !ok {
		chart.Dependencies = []Chart{}
	}
	for _, d := range deps {
		dep, ok := d.(map[interface{}]interface{})
		if !ok {
			continue
		}
		for k, v := range dep {
			name := k.(string)
			extra := v.(map[interface{}]interface{})
			r, ok := extra["repo"].(string)
			if !ok {
				r = ""
			}
			v, ok := extra["version"].(string)
			if !ok {
				v = ""
			}
			chart.Dependencies = append(chart.Dependencies, Chart{Name: name, Repo: r, Version: v})
		}
	}
}

func getPrereqCharts() (retCharts []Chart, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to parse prereq charts in config file")
			retCharts = nil
		}
	}()

	prereqList := viper.Get("prereqCharts").([]interface{})
	for _, chart := range prereqList {
		for name, vals := range chart.(map[interface{}]interface{}) {
			vals := vals.(map[interface{}]interface{})
			nextChart := Chart{
				Name: name.(string),
			}

			r, ok := vals["repo"].(string)
			if ok {
				nextChart.Repo = r
			}

			v, ok := vals["version"].(string)
			if ok {
				nextChart.Version = v
			}

			c, ok := vals["pathToChart"].(string)
			if ok {
				nextChart.PathToChart = c
			}

			p, ok := vals["packaged"].(bool)
			if ok {
				nextChart.Packaged = p
			} else {
				nextChart.Packaged = true
			}

			files := []string{}
			overrideFiles, ok := vals["overrideValueFiles"].([]interface{})
			if ok {
				for _, f := range overrideFiles {
					files = append(files, f.(string))
				}
			}
			nextChart.OverrideValueFiles = files

			retCharts = append(retCharts, nextChart)
		}
	}

	return retCharts, err
}
