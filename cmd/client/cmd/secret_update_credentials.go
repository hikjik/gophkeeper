package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var updateCredentialsSecretCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Update credentials secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
		}

		login, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Fatal().Msgf("Error reading login: %v", err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal().Msgf("Error reading password: %v", err)
		}

		credentials := models.Credentials{
			Login:    login,
			Password: password,
		}

		content, err := encryptSecret(credentials)
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
	updateSecretCmd.AddCommand(updateCredentialsSecretCmd)

	updateCredentialsSecretCmd.Flags().String("name", "", "Secret name")
	if err := updateCredentialsSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	updateCredentialsSecretCmd.Flags().String("login", "", "Login")
	if err := updateCredentialsSecretCmd.MarkFlagRequired("login"); err != nil {
		log.Error().Err(err)
	}
	updateCredentialsSecretCmd.Flags().String("password", "", "Password")
	if err := updateCredentialsSecretCmd.MarkFlagRequired("password"); err != nil {
		log.Error().Err(err)
	}
}
