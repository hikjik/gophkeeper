package cmd

import (
	"github.com/spf13/cobra"
)

var createSecretCmd = &cobra.Command{
	Use:   "create",
	Short: "Create secret",
}

func init() {
	secretCmd.AddCommand(createSecretCmd)
}
