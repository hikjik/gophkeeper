package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logins a user in the gophkeeper service",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read email")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read password")
		}

		resp, err := client.SignIn(context.Background(), &proto.SignInRequest{Email: email, Password: password})
		if err != nil {
			log.Error().Err(err).Msg("Failed to login")
			return
		}

		log.Info().Msgf("Access Token: %s", resp.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "User Email")
	if err := loginCmd.MarkFlagRequired("email"); err != nil {
		log.Error().Err(err)
	}
	loginCmd.Flags().StringP("password", "p", "", "User password")
	if err := loginCmd.MarkFlagRequired("password"); err != nil {
		log.Error().Err(err)
	}
}
