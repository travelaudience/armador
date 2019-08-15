package cmd

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/travelaudience/armador/internal/cluster"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/helm"
	"github.com/travelaudience/armador/internal/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [namespace/env]",
	Short: "Tear down an existing dev env",
	Run: func(cmd *cobra.Command, args []string) {
		nsToRemove := args[0]
		logger.GetLogger().Debugf("the ns to remove is %s", nsToRemove)
		delete(nsToRemove)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("delete requires a namespace/env to remove")
		}
		if len(args) > 1 {
			return errors.New("delete can only remove one namespace/env")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func delete(nsToRemove string) error {
	logger := logger.GetLogger()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	cmd := commands.Cmd{Ctx: ctx}

	logger.Infof("Deleting %s...", nsToRemove)

	//  connect to cluster
	clusterConfig := viper.Get("cluster").(cluster.ClusterConfig)
	err := cluster.ClusterConnect(cmd, clusterConfig)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if !cluster.NamespaceExists(cmd, nsToRemove) {
		logger.Warnf("%s is not available in this cluster", nsToRemove)
		return nil
	}
	// TODO: CONFIG: create global settings for protected namespaces
	// 		 these should be built into the app, but easy to change if someone forks the repo
	if strings.Contains(nsToRemove, "production") || strings.Contains(nsToRemove, "staging") {
		logger.Warnf("%s is one of the protected namespaces and must be managed manually", nsToRemove)
		return nil
	}

	helm.PurgeInstalls(cmd, nsToRemove)
	cluster.DeleteNamespace(cmd, nsToRemove)
	return nil
}
