package cmd

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/travelaudience/armador/internal/armador"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/cluster"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot [namespace/env]",
	Short: "Create a snapshot of a running env, using helm charts/values save to config file.",
	Long: `Create a snapshot of a running env, using helm charts/values save to config file.

Given the details of a running environment, this command will establish what helm
charts are currently installed there, and what values were used during their release.
This data will be saved to a configuration files, and will allow for the same setup to
be created in a new environment (using 'armador recreate'). It also can be used similar
to 'heptio/velero' for creating backups (however it does not handle persistent data).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("`snapshot` requires a namespace/env to copy")
		}
		if len(args) > 1 {
			return errors.New("`snapshot` can only manage one namespace/env at a time")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		nsToSave := args[0]
		overridePathToSave, err := cmd.Flags().GetString("path")
		if err != nil {
			logger.GetLogger().Warnf("Problem reading parameters: %s", err)
			return
		}
		logger.GetLogger().Debugf("the ns to copy is %s", nsToSave)
		snapshotCreation(nsToSave, overridePathToSave)
	},
}

func init() {
	snapshotCmd.Flags().StringP("path", "p", "", "path to the snapshot configuration (defaults to $ARMADOR_HOME/snapshot/[namespace])")
	rootCmd.AddCommand(snapshotCmd)
}

func snapshotCreation(nsToSave, pathToSave string) {
	logger := logger.GetLogger()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	cmd := commands.Cmd{Ctx: ctx}
	logger.Infof("Creating snapshot of %s...", nsToSave)

	//  connect to cluster
	clusterConfig := viper.Get("cluster").(cluster.ClusterConfig)
	err := cluster.ClusterConnect(cmd, clusterConfig)
	if err != nil {
		logger.Error(err)
		return
	}

	if !cluster.NamespaceExists(cmd, nsToSave) {
		logger.Warnf("%s is not available in this cluster: %v", nsToSave, clusterConfig)
		return
	}

	if pathToSave == "" {
		dirs, err := commands.CreateDirs()
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debugf("snapshot parent dir: %s ", dirs.Tmp)
		pathToSave = filepath.Join(dirs.Snapshot, nsToSave)
	}
	// make sure `pathToSave` exists (and sub folder /overrides for the values)
	err = commands.CheckDir(pathToSave)
	if err != nil {
		logger.Error(err)
		return
	}
	valuesPath := filepath.Join(pathToSave, "overrides")
	err = commands.CheckDir(valuesPath)
	if err != nil {
		logger.Error(err)
		return
	}

	err = armador.CreateSnapshot(cmd, nsToSave, pathToSave, valuesPath)
	if err != nil {
		logger.Error(err)
		return
	}
}
