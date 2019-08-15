package cmd

import (
	"github.com/spf13/cobra"
	"github.com/travelaudience/armador/internal/commands"
	"github.com/travelaudience/armador/internal/logger"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove cache and temp files",
	Run: func(cmd *cobra.Command, args []string) {
		clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func clean() {
	logger := logger.GetLogger()
	defer logger.Sync()

	logger.Info("Cleaning cache...")

	err := commands.CleanDirs()
	if err != nil {
		logger.Errorf("Clearing directories failed: %s", err)
	}
	logger.Debugf("temp & cache dir cleared")
}
