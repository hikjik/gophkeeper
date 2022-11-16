package server

import (
	"context"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/config"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/services"
)

// Server сервер gRPC с сервисами аутентификации и хранения пользовательских данных
type Server struct {
	AuthService   *services.AuthService
	SecretService *services.SecretService
	Address       string
}

// New создает новый Server с указанными настройками
func New(cfg config.Config) *Server {
	authService := services.NewAuthService(cfg)
	secretService := services.NewSecretService(cfg)
	return &Server{
		AuthService:   authService,
		SecretService: secretService,
		Address:       cfg.GRPC.Address,
	}
}

// Run функция запуска gRPC сервера с сервисами аутентификации и хранения пользовательских данных
func (s *Server) Run(ctx context.Context) {
	interceptor := interceptors.NewAuthInterceptor(s.AuthService.TokenManager)

	services.NewServer(
		s.Address,
		services.WithServices(s.AuthService, s.SecretService),
		services.WithUnaryInterceptors(interceptor.Unary()),
		services.WithStreamInterceptors(interceptor.Stream()),
	).Run(ctx)
}
