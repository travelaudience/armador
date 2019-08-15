package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Manage your local helm install",
	Long: `This command doesn't do anything on it's own
But can be used with the additional [commands] to manage your local helm setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("an additonal command is required")
	},
}

func init() {
	rootCmd.AddCommand(helmCmd)
	helmCmd.AddCommand(newHelmCheckCmd())
}
