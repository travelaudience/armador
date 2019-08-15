package helm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/travelaudience/armador/internal/commands"
)

func Install(cmd commands.Command, name, chartPath, namespace, overridePath string, overrides, setValues []string) (string, error) {
	// run this cmds in the `chartPath` dir
	installCmds := []string{"helm", "upgrade", "--install", name + "-" + namespace, ".", "--namespace", namespace}
	if len(overrides) > 0 {
		installCmds = append(installCmds, addFlag(overrides, "-f")...)
	}
	additionalOverride := filepath.Join(overridePath, name+".yaml")
	if _, err := os.Stat(additionalOverride); err == nil {
		installCmds = append(installCmds, "-f", additionalOverride)
	}
	if len(setValues) > 0 {
		installCmds = append(installCmds, addFlag(setValues, "--set")...)
	}
	return cmd.ExecInDirUnparsed("helm-install", chartPath, installCmds...)
}

func Diff(cmd commands.Command, name, chartPath, namespace, overridePath string, overrides, setValues []string) (string, error) {
	// run this cmds in the `chartPath` dir
	dryRunCmds := []string{"helm", "diff", "upgrade", name + "-" + namespace, ".", "--allow-unreleased"}
	if len(overrides) > 0 {
		dryRunCmds = append(dryRunCmds, addFlag(overrides, "-f")...)
	}
	additionalOverride := filepath.Join(overridePath, name+".yaml")
	if _, err := os.Stat(additionalOverride); err == nil {
		dryRunCmds = append(dryRunCmds, "-f", additionalOverride)
	}
	if len(setValues) > 0 {
		dryRunCmds = append(dryRunCmds, addFlag(setValues, "--set")...)
	}
	return cmd.ExecInDirUnparsed("helm-diff", chartPath, dryRunCmds...)
}

func InstallDirect(cmd commands.Command, chartName, namespace, chartPath, overridePath string) {
	dryRunCmds := []string{"helm", "upgrade", "--install", chartName + "-" + namespace, ".", "--namespace", namespace}
	additionalOverride := filepath.Join(overridePath, chartName+".yaml")
	dryRunCmds = append(dryRunCmds, "-f", additionalOverride)
	cmd.ExecInDir("helm-install", chartPath, dryRunCmds...)
}

func DiffDirect(cmd commands.Command, chartName, namespace, chartPath, overridePath string) {
	dryRunCmds := []string{"helm", "diff", "upgrade", chartName + "-" + namespace, ".", "--allow-unreleased"}
	additionalOverride := filepath.Join(overridePath, chartName+".yaml")
	dryRunCmds = append(dryRunCmds, "-f", additionalOverride)
	cmd.ExecInDir("helm-diff", chartPath, dryRunCmds...)
}

func Fetch(cmd commands.Command, chart, repo, version, holdDir, extractDir, cacheDir string) (string, error) {
	tarFileName := ""
	fetchCmds := []string{"helm", "fetch", "-d", holdDir, repo + "/" + chart}
	if version != "" {
		fetchCmds = append(fetchCmds, "--version", version)
	}
	// TODO: maybe make sure the hold dir is empty first
	_, err := cmd.Exec("helm-fetch", fetchCmds...)
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(holdDir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), chart) {
			// TODO: remove the file
			continue
		}
		tarFileName = file.Name()
	}
	if version == "" {
		version = getVersionFromFilename(tarFileName, chart)
	}
	dir := filepath.Join(cacheDir, chart, version)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Failed to create %s dir: %s", dir, err)
	}
	// move the hold file to the cache dir
	err = os.Rename(filepath.Join(holdDir, tarFileName), filepath.Join(dir, tarFileName))
	if err != nil {
		return "", fmt.Errorf("Failed to move tar file %s to dir: %s: %s", tarFileName, dir, err)
	}
	// extract the cached tar file to the tmp/extracted dir
	err = commands.Extract(dir, extractDir, tarFileName)
	if err != nil {
		return "", fmt.Errorf("Extraction of %s failed: %s", chart, err)
	}
	return filepath.Join(extractDir, chart), nil
}

func ListCharts(cmd commands.Command, namespace string) ([]string, error) {
	// TODO: CONFIG: make this an optional value
	maxHelmList := "80"
	return cmd.Exec("helm-list", "helm", "list", "-qadr", "-m", maxHelmList, "--namespace", namespace)
}

func getVersionFromFilename(filename, chartName string) string {
	filename = strings.Replace(filename, chartName, "", 1)
	lsIn := strings.LastIndex(filename, ".")
	return filename[1:lsIn]
}

func GetChartVersion(cmd commands.Command, chartName string) (string, error) {
	return cmd.ExecUnparsed("helm-version", "helm", "get", chartName, "--template", "{{.Release.Chart.Metadata.Version}}")
}

func CollectValues(cmd commands.Command, chartName string) (string, error) {
	return cmd.ExecUnparsed("helm-get", "helm", "get", "values", chartName)
}

func RepoUpdate(cmd commands.Command) {
	cmd.Exec("helm-repo-update", "helm", "repo", "update")
}

func addFlag(s []string, key string) []string {
	if len(s) < 1 {
		return []string{}
	}
	keyBuf := " " + key + " "
	flattened := keyBuf + strings.Join(s, keyBuf)
	withF := strings.Split(flattened, " ")
	var r []string
	for _, str := range withF {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func PurgeInstalls(cmd commands.Command, namespace string) error {
	var wg sync.WaitGroup

	charts, err := ListCharts(cmd, namespace)
	if err != nil {
		return err
	}

	wg.Add(len(charts))
	for _, chart := range charts {
		go func(cmd commands.Command, chartName string) {
			defer wg.Done()
			cmd.Exec("helm-purge", "helm", "delete", "--purge", chartName)
		}(cmd, chart)
	}
	wg.Wait()
	return nil
}

func CheckPlugin(cmd commands.Command, pluginName string) error {
	res, err := cmd.Exec("helm-plugin", "helm", "plugin", "list", pluginName)
	if err != nil {
		return err
	}
	for _, n := range res {
		if strings.HasPrefix(n, pluginName+" ") {
			return nil
		}
	}
	return fmt.Errorf("Helm plugin: %s not installed. /nTry `armardor helm check` to compare versions", pluginName)
}

func DepUpdate(cmd commands.Command, chartPath string) ([]string, error) {
	return cmd.ExecInDir("helm-dep-up", chartPath, "helm", "dep", "update")
}
