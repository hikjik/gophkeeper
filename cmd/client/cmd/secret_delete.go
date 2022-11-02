package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var deleteSecretCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
		}

		resp, err := secretClient.DeleteSecret(
			context.Background(), &pb.DeleteSecretRequest{Name: name})
		if err != nil {
			log.Fatal().Msgf("Failed to delete secret: %v", err)
			return
		}

		fmt.Printf("Secret %s deleted successfully\n", resp.GetName())
	},
}

func init() {
	secretCmd.AddCommand(deleteSecretCmd)

	deleteSecretCmd.Flags().String("name", "", "Secret name")
	if err := deleteSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
}
