package cmd

import (
	"github.com/spf13/cobra"
)

var updateSecretCmd = &cobra.Command{
	Use:   "update",
	Short: "Update secret",
}

func init() {
	secretCmd.AddCommand(updateSecretCmd)
}
