package cmd

import (
	"context"
	"fmt"

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

		resp, err := authClient.VerifyToken(
			context.Background(),
			&proto.VerifyTokenRequest{AccessToken: token},
		)
		if err != nil {
			fmt.Printf("Token is invalid: %v\n", err)
			return
		}

		fmt.Printf("Token is valid, UserID: %d\n", resp.UserId)
	},
}

func init() {
	authCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringP("token", "t", "", "Access token")
	if err := verifyCmd.MarkFlagRequired("token"); err != nil {
		log.Error().Err(err)
	}
}
