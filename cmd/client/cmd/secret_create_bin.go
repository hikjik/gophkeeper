package cmd

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var createBinSecretCmd = &cobra.Command{
	Use:   "bin",
	Short: "Create bin secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
			return
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal().Msgf("Error reading file name: %v", err)
			return
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal().Msgf("Error reading binary file: %v", err)
			return
		}

		bin := models.Bin{
			Data: data,
		}

		content, err := encryptSecret(bin)
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
	createSecretCmd.AddCommand(createBinSecretCmd)

	createBinSecretCmd.Flags().String("name", "", "Secret name")
	if err := createBinSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	createBinSecretCmd.Flags().StringP("file", "f", "", "Binary file")
	if err := createBinSecretCmd.MarkFlagRequired("file"); err != nil {
		log.Error().Err(err)
	}
}
