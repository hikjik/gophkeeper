package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies access token",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read token")
		}

		resp, err := client.VerifyToken(context.Background(), &proto.VerifyTokenRequest{AccessToken: token})
		if err != nil {
			log.Error().Err(err).Msg("Failed to verify token")
			return
		}

		log.Info().Msgf("Token is valid, user id: %d", resp.UserId)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringP("token", "t", "", "Access token")
	if err := verifyCmd.MarkFlagRequired("token"); err != nil {
		log.Error().Err(err)
	}
}
