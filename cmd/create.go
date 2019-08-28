package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/travelaudience/armador/internal/armador"
	"github.com/travelaudience/armador/internal/cluster"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create [namespace/env]",
	Short: "Create a dev env",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("`create` requires a namespace/env to create")
		}
		if len(args) > 1 {
			return errors.New("`create` can only manage one namespace/env at a time")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		namespace := args[0]
		projectDir, err := cmd.Flags().GetString("projectDir")
		if err != nil {
			return err
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		rawValues, err := cmd.Flags().GetStringArray("set")
		if err != nil {
			return err
		}
		logger.GetLogger().Debugf("the flags are %s, %s\n", projectDir, namespace)
		logger.GetLogger().Debugf("raw value overrides %v\n", rawValues)
		logger.GetLogger().Debugf("dry run is %v\n", dryRun)
		createOrUpdate(projectDir, namespace, dryRun, true, rawValues)
		return nil
	},
}

// The purpose of the `update` cmd is to give users the awarness that their env already exists or not.
//   functionality is in the end the same because it uses `helm upgrade --install`
var updateCmd = &cobra.Command{
	Use:     "update [namespace/env]",
	Short:   "Update a dev env",
	Aliases: []string{"update"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("`update` requires a namespace/env to update")
		}
		if len(args) > 1 {
			return errors.New("`update` can only manage one namespace/env at a time")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		namespace := args[0]
		projectDir, err := cmd.Flags().GetString("projectDir")
		if err != nil {
			return err
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		rawValues, err := cmd.Flags().GetStringArray("set")
		if err != nil {
			return err
		}
		logger.GetLogger().Debugf("the flags are %s, %s\n", projectDir, namespace)
		logger.GetLogger().Debugf("raw value overrides %v\n", rawValues)
		logger.GetLogger().Debugf("dry run is %v\n", dryRun)
		createOrUpdate(projectDir, namespace, dryRun, false, rawValues)
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("projectDir", "p", ".", "path to the project that will be installed")
	createCmd.Flags().BoolP("dry-run", "d", false, "show the changes helm will apply without making them")
	createCmd.Flags().StringArray("set", []string{}, "set values on the command line (can specify multiple: --set app.key=val1 --set app2.key.v2=val2)")
	rootCmd.AddCommand(createCmd)
	// Set the same for update as well
	updateCmd.Flags().StringP("projectDir", "p", ".", "path to the project that will be installed")
	updateCmd.Flags().BoolP("dry-run", "d", false, "show the changes helm will apply without making them")
	updateCmd.Flags().StringArray("set", []string{}, "set values on the command line (can specify multiple: --set app.key=val1 --set app2.key.v2=val2)")
	rootCmd.AddCommand(updateCmd)
}

func createOrUpdate(projectDir, namespace string, dryRun, create bool, rawValues []string) {
	logger := logger.GetLogger()
	defer logger.Sync()

	// TODO: set the minutes as a config option
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := commands.Cmd{Ctx: ctx}
	if create {
		logger.Info("Creating...")
	} else {
		logger.Info("Updating...")
	}

	dirs, err := commands.CreateDirs()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debugf("temp & cache dir: %s & %s ", dirs.Tmp, dirs.Cache)

	//  connect to cluster
	clusterConfig := viper.Get("cluster").(cluster.ClusterConfig)
	err = cluster.ClusterConnect(cmd, clusterConfig)
	if err != nil {
		logger.Warn(err)
		return
	}

	// validate the namespace and if ok, create it already
	err = cluster.CheckName(namespace)
	if err != nil {
		logger.Warn(err)
		return
	}
	nameExists := cluster.NamespaceExists(cmd, namespace)
	if create && nameExists {
		logger.Warnf("%s already exists. Either `armador delete` it, or use `armador update`.", namespace)
		return
	}
	if !create && !nameExists {
		logger.Warnf("%s does not exist. Use `armador create` instead.", namespace)
		return
	}
	if create && !dryRun {
		_, err = cluster.CreateNamespace(cmd, namespace)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	err = armador.Create(cmd, projectDir, namespace, dirs, dryRun, rawValues)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Creation is complete, enjoy using %s", namespace)

	err = cluster.SetContextNamespace(cmd, namespace)
	if err != nil {
		logger.Warn("Unable to set namespace in context.")
	}
	if dryRun {
		logger.Infof("Dry-run is complete, this is what will be deployed.")
	} else {
		logger.Infof("Creation is complete, enjoy using %s", namespace)
	}

}
