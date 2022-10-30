package server

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/config"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/hasher"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/hasher/hmac"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage/pg"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/token"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/token/jwt"
)

// Server реализация proto.GophKeeperServer
type Server struct {
	proto.UnimplementedGophKeeperServer

	Storage      storage.Storage
	TokenManager token.Manager
	Hasher       hasher.Hasher
	Address      string
}

var _ proto.GophKeeperServer = (*Server)(nil)

// New создает новый Server
func New(cfg config.Config) *Server {
	dbStorage, err := pg.New(cfg.DB.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create storage")
	}

	tokenManager, err := jwt.New(cfg.Auth.Key, cfg.Auth.ExpirationTime)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create token manager")
	}

	hmacHasher, err := hmac.New(cfg.Hash.Key)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create hasher computer")
	}

	return &Server{
		Storage:      dbStorage,
		TokenManager: tokenManager,
		Hasher:       hmacHasher,
		Address:      cfg.GRPC.Address,
	}
}

// Run функция запуска сервера
func (s *Server) Run(ctx context.Context) {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start grpc server")
	}

	interceptor := interceptors.NewAuthInterceptor(s.TokenManager)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	proto.RegisterGophKeeperServer(server, s)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	if err = server.Serve(listen); err != nil {
		log.Error().Err(err).Msg("Error on grpc server Serve")
	}
}
