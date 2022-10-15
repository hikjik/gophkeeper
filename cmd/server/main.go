package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/go-developer-ya-practicum/gophkeeper/cmd/server/cmd"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/greeting"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if err := greeting.PrintBuildInfo(os.Stdout, buildVersion, buildDate, buildCommit); err != nil {
		log.Warn().Err(err).Msg("Failed to print build info")
	}

	cmd.Execute()
}
