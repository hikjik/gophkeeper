package cmd

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/client/config"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/version"
)

var (
	cfgFile  string
	defaults = map[string]interface{}{
		"grpc.address": "127.0.0.1:8081",
	}

	client pb.GophKeeperClient
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gophkeeper-cli",
	Short:   "GophKeeper client",
	Version: version.Info(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal().Err(err).Msg("Failed to load client config")
		}

		connection, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal().Err(err)
		}

		client = pb.NewGophKeeperClient(connection)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "", "Client config filepath")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/gophkeeper")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")

		viper.SetConfigName("client-config")
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal().Err(err).Msg("Failed to load server config")
		}
	} else {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		key := strings.ReplaceAll(flag.Name, "-", ".")
		if err := viper.BindPFlag(key, flag); err != nil {
			log.Fatal().Err(err).Msg("Failed to bind flag")
		}
	})
}
