package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/helm"
	"github.com/travelaudience/armador/internal/logger"
)

func newHelmCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Output the state of current helm setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck()
		},
	}

	return cmd
}

func runCheck() error {
	logger := logger.GetLogger()
	defer logger.Sync()

	// TODO: set the minutes as a config option
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := commands.Cmd{Ctx: ctx}

	// check helm version
	ver, err := cmd.ExecUnparsed("helm-version", "helm", "version")
	if err != nil {
		logger.Warnf("Unable to obtain helm version info: %s", err)
		return err
	}
	logger.Infof("Helm Version Info: \n%s", ver)

	// check helm plugins
	plugins, err := cmd.ExecUnparsed("helm-plugins", "helm", "plugin", "list")
	if err != nil {
		logger.Warnf("Unable to obtain helm plugins: %s", err)
		return err
	}
	logger.Infof("Helm plugins: \n%s", plugins)

	// expected plugins:
	pluginList := []string{"diff"}
	for _, p := range pluginList {
		err = helm.CheckPlugin(cmd, p)
		if err != nil {
			logger.Warnf("Plugin: %s not found, but may be required.", p)
		}
	}

	// list helm repos
	repos, err := cmd.ExecUnparsed("helm-repos", "helm", "repo", "list")
	if err != nil {
		logger.Warnf("Unable to obtain helm repos: %s", err)
		return err
	}
	logger.Infof("Helm repos: \n%s", repos)

	return nil
}
