package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/models"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/cipher"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/cipher/aes/gcm"
)

var (
	secretClient pb.SecretServiceClient
	blockCipher  cipher.BlockCipher
)

func encryptSecret(s models.Secret) ([]byte, error) {
	encoded, err := models.EncodeSecret(s)
	if err != nil {
		return nil, err
	}
	return blockCipher.Encrypt(encoded)
}

func decryptSecret(b []byte) (models.Secret, error) {
	encoded, err := blockCipher.Decrypt(b)
	if err != nil {
		return nil, err
	}
	return models.DecodeSecret(encoded)
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage user private data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		accessToken := viper.GetString("token")
		if accessToken == "" {
			log.Fatal().Msg("Empty access token")
		}
		interceptor := interceptors.NewAuthInterceptor(viper.GetString("token"))

		connection, err := grpc.Dial(
			viper.GetString("grpc.address"),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(interceptor.Unary()),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create client connection")
		}

		secretClient = pb.NewSecretServiceClient(connection)
		cipher, err := gcm.New(viper.GetString("encryption.key"))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create cipher")
		}
		blockCipher = cipher
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
