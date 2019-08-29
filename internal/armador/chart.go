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
	vip, err := ReadFileToViper("armador", chart.ChartPath) // TODO: CONFIG: the config name should be set somehow
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

func digestFields(fields map[interface{}]interface{}) (retChart Chart) {
	retChart = Chart{}
	for name, vals := range fields {
		retChart.Name = name.(string)
		vals := vals.(map[interface{}]interface{})
		r, ok := vals["repo"].(string)
		if ok {
			retChart.Repo = r
		}
		v, ok := vals["version"].(string)
		if ok {
			retChart.Version = v
		}
		c, ok := vals["pathToChart"].(string)
		if ok {
			retChart.PathToChart = c
		}

		p, ok := vals["packaged"].(bool)
		if ok {
			retChart.Packaged = p
		} else {
			retChart.Packaged = true
		}

		var files []string
		overrideFiles, ok := vals["overrideValueFiles"].([]interface{})
		if ok {
			for _, f := range overrideFiles {
				files = append(files, f.(string))
			}
		}
		retChart.OverrideValueFiles = files
	}
	return retChart
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

	deps, ok := vip.Get("dependencies").([]interface{})
	if !ok {
		chart.Dependencies = []Chart{}
	}
	for _, d := range deps {
		dep, _ := d.(map[interface{}]interface{})
		chart.Dependencies = append(chart.Dependencies, digestFields(dep))
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
		preq, _ := chart.(map[interface{}]interface{})
		retCharts = append(retCharts, digestFields(preq))
	}

	return retCharts, err
}
