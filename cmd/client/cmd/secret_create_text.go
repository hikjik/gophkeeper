package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var createTextSecretCmd = &cobra.Command{
	Use:   "text",
	Short: "Create text secret",
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

		resp, err := secretClient.CreateSecret(context.Background(), &pb.CreateSecretRequest{
			Name:    name,
			Content: content,
		})
		if err != nil {
			log.Fatal().Msgf("Failed to create secret: %v", err)
			return
		}

		fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	createSecretCmd.AddCommand(createTextSecretCmd)

	createTextSecretCmd.Flags().String("name", "", "Secret name")
	if err := createTextSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	createTextSecretCmd.Flags().String("data", "", "Text data")
	if err := createTextSecretCmd.MarkFlagRequired("data"); err != nil {
		log.Error().Err(err)
	}
}
