package server

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/hasher"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/hasher/hmac"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/storage"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/storage/pg"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/token"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/token/jwt"
)

// Server реализация proto.GophKeeperServer
type Server struct {
	proto.UnimplementedGophKeeperServer

	Storage        storage.Storage
	TokenGenerator token.Generator
	Hasher         hasher.Hasher
	Address        string
}

var _ proto.GophKeeperServer = (*Server)(nil)

// New создает новый Server
func New(cfg Config) *Server {
	dbStorage, err := pg.New(cfg.DB.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create storage")
	}

	tokenGenerator, err := jwt.New(cfg.Auth.Key, cfg.Auth.ExpirationTime)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create token generator")
	}

	hmacHasher, err := hmac.New(cfg.Hash.Key)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create hasher computer")
	}

	return &Server{
		Storage:        dbStorage,
		TokenGenerator: tokenGenerator,
		Hasher:         hmacHasher,
		Address:        cfg.GRPC.Address,
	}
}

// Run функция запуска сервера
func (s *Server) Run(ctx context.Context) {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start grpc server")
	}
	server := grpc.NewServer()
	proto.RegisterGophKeeperServer(server, s)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	if err = server.Serve(listen); err != nil {
		log.Error().Err(err).Msg("Error on grpc server Serve")
	}
}
