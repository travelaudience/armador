package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/armador"
	"github.com/travelaudience/armador/internal/cluster"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

// recreateCmd represents the recreate command
var recreateCmd = &cobra.Command{
	Use:   "recreate [namespace/env]",
	Short: "Use a cache file to re-create an env with the same settings",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("recreate requires a namespace/env (either the existing one, or a new one)")
		}
		if len(args) > 1 {
			return errors.New("recreate can only manage one namespace/env at a time")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]
		dryRun, err := cmd.Flags().GetBool("dryRun")
		if err != nil {
			return
		}
		recreate(namespace, dryRun)
	},
}

func init() {
	recreateCmd.Flags().BoolP("dryRun", "d", false, "show the changes helm will apply without making them")
	rootCmd.AddCommand(recreateCmd)

}

func recreate(namespace string, dryRun bool) {
	logger := logger.GetLogger()
	defer logger.Sync()

	// TODO: set the minutes as a config option
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := commands.Cmd{Ctx: ctx}
	logger.Info("Re-creating...")

	dirs, err := commands.CreateDirs()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debugf("temp & cache dir: %s & %s ", dirs.Tmp, dirs.Cache)

	// TODO: CONFIG: set this in a global config
	cacheFile := filepath.Join(dirs.Cache.Root, "envCache.yaml")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		logger.Warnf("Cache file: %s does not exist. Unable to recreate from nothing... Try using `create`", cacheFile)
		return
	}

	//  connect to cluster
	clusterConfig := viper.Get("cluster").(cluster.ClusterConfig)
	err = cluster.ClusterConnect(cmd, clusterConfig)
	if err != nil {
		logger.Error(err)
		return
	}

	// TODO: does it make sense to check if namespace exists, and warn/abort

	err = armador.Recreate(cmd, cacheFile, namespace, dirs, dryRun)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Re-creation is complete, enjoy using %s", namespace)
}
