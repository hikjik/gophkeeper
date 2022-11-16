package cmd

import (
	"context"
	"fmt"

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

		resp, err := authClient.SignIn(
			context.Background(),
			&proto.SignInRequest{Email: email, Password: password},
		)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return
		}

		if err := tokenStorage.Save(resp.AccessToken); err != nil {
			log.Fatal().Err(err).Msg("Failed to store access token")
		}
		fmt.Printf("Access Token: %s\n", resp.AccessToken)
	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "User Email")
	if err := loginCmd.MarkFlagRequired("email"); err != nil {
		log.Error().Err(err)
	}
	loginCmd.Flags().StringP("password", "p", "", "User password")
	if err := loginCmd.MarkFlagRequired("password"); err != nil {
		log.Error().Err(err)
	}
}
