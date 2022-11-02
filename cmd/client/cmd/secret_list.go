package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var listSecretCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := secretClient.ListSecrets(context.Background(), &pb.ListSecretsRequest{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list secret")
		}

		for _, info := range resp.GetSecrets() {
			secret, err := decryptSecret(info.GetContent())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to decrypt secret")
			}

			fmt.Printf("%s\n", secret)
		}
	},
}

func init() {
	secretCmd.AddCommand(listSecretCmd)
}
