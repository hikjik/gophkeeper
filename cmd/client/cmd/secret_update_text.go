package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var updateTextSecretCmd = &cobra.Command{
	Use:   "text",
	Short: "Update text secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
			return
		}

		data, err := cmd.Flags().GetString("data")
		if err != nil {
			log.Fatal().Msgf("Error reading text data: %v", err)
			return
		}

		text := models.Text{
			Data: data,
		}

		content, err := encryptSecret(text)
		if err != nil {
			log.Fatal().Msgf("Failed to encrypt secret: %v", err)
			return
		}

		resp, err := secretClient.UpdateSecret(context.Background(), &pb.UpdateSecretRequest{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Fatal().Msgf("Failed to update secret: %v", err)
			return
		}

		fmt.Printf("Secret %s version %v updated successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	updateSecretCmd.AddCommand(updateTextSecretCmd)

	updateTextSecretCmd.Flags().String("name", "", "Secret name")
	if err := updateTextSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	updateTextSecretCmd.Flags().String("data", "", "Text data")
	if err := updateTextSecretCmd.MarkFlagRequired("data"); err != nil {
		log.Error().Err(err)
	}
}
