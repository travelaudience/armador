package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/travelaudience/armador/internal/cluster"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the environments currently running",
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list() {
	logger := logger.GetLogger()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	cmd := commands.Cmd{Ctx: ctx}

	logger.Debug("List armador environments ...")

	//  connect to cluster
	clusterConfig := viper.Get("cluster").(cluster.ClusterConfig)
	err := cluster.ClusterConnect(cmd, clusterConfig)
	if err != nil {
		logger.Error(err)
		return
	}

	envs := cluster.ListEnvironments(cmd)
	logger.Infof("Potential Armador environments:\n%s", strings.Join(envs[:], "\n"))
}
