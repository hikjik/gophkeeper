package cmd

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/config"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "GophKeeper server daemon",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(
			context.Background(), syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal().Err(err).Msg("Failed to load server config")
		}
		log.Info().Msgf("Start grpc server: %s", cfg.GRPC.Address)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.NewDaemon(cfg).Run(ctx)
		}()
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
