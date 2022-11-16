package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var createCredentialsSecretCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Create credentials secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
			return
		}

		login, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Fatal().Msgf("Error reading login: %v", err)
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal().Msgf("Error reading password: %v", err)
			return
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
	createSecretCmd.AddCommand(createCredentialsSecretCmd)

	createCredentialsSecretCmd.Flags().String("name", "", "Secret name")
	if err := createCredentialsSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	createCredentialsSecretCmd.Flags().String("login", "", "Login")
	if err := createCredentialsSecretCmd.MarkFlagRequired("login"); err != nil {
		log.Error().Err(err)
	}
	createCredentialsSecretCmd.Flags().String("password", "", "Password")
	if err := createCredentialsSecretCmd.MarkFlagRequired("password"); err != nil {
		log.Error().Err(err)
	}
}
