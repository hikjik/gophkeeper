package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a user in the gophkeeper service",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read email")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read password")
		}

		resp, err := client.SignUp(context.Background(), &proto.SignUpRequest{Email: email, Password: password})
		if err != nil {
			log.Error().Err(err).Msg("Failed to register")
			return
		}

		log.Info().Msgf("Access Token: %s", resp.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("email", "e", "", "User Email")
	if err := registerCmd.MarkFlagRequired("email"); err != nil {
		log.Error().Err(err)
	}

	registerCmd.Flags().StringP("password", "p", "", "User password")
	if err := registerCmd.MarkFlagRequired("password"); err != nil {
		log.Error().Err(err)
	}
}
