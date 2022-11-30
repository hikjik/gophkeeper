package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var getSecretCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
		}

		resp, err := secretClient.GetSecret(context.Background(), &pb.GetSecretRequest{
			Name: name,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get secret")
		}

		secret, err := decryptSecret(resp.GetContent())
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to decrypt secret")
		}

		fmt.Printf("%s\n", secret)
	},
}

func init() {
	secretCmd.AddCommand(getSecretCmd)

	getSecretCmd.Flags().String("name", "", "Secret name")
	if err := getSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
}
