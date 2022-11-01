package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
)

var (
	authClient pb.AuthServiceClient
)

// authCmd represents the solution command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		connection, err := grpc.Dial(
			viper.GetString("grpc.address"),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create client connection")
		}

		authClient = pb.NewAuthServiceClient(connection)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
