package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var createCardSecretCmd = &cobra.Command{
	Use:   "card",
	Short: "Create card secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
			return
		}

		number, err := cmd.Flags().GetString("number")
		if err != nil {
			log.Fatal().Msgf("Error reading card number: %v", err)
			return
		}

		date, err := cmd.Flags().GetString("date")
		if err != nil {
			log.Fatal().Msgf("Error reading card expiry date: %v", err)
			return
		}

		code, err := cmd.Flags().GetString("code")
		if err != nil {
			log.Fatal().Msgf("Error reading card security code: %v", err)
			return
		}

		holder, err := cmd.Flags().GetString("holder")
		if err != nil {
			log.Fatal().Msgf("Error reading card holder: %v", err)
			return
		}

		card := models.Card{
			Number:       number,
			ExpiryDate:   date,
			SecurityCode: code,
			Holder:       holder,
		}

		content, err := encryptSecret(card)
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
	createSecretCmd.AddCommand(createCardSecretCmd)

	createCardSecretCmd.Flags().String("name", "", "Secret name")
	if err := createCardSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	createCardSecretCmd.Flags().String("number", "", "Card number")
	if err := createCardSecretCmd.MarkFlagRequired("number"); err != nil {
		log.Error().Err(err)
	}
	createCardSecretCmd.Flags().String("date", "", "Card expiry date")
	if err := createCardSecretCmd.MarkFlagRequired("date"); err != nil {
		log.Error().Err(err)
	}
	createCardSecretCmd.Flags().String("code", "", "Card security code")
	if err := createCardSecretCmd.MarkFlagRequired("code"); err != nil {
		log.Error().Err(err)
	}
	createCardSecretCmd.Flags().String("holder", "", "Card holder")
	if err := createCardSecretCmd.MarkFlagRequired("holder"); err != nil {
		log.Error().Err(err)
	}
}
