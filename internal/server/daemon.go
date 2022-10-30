package server

import (
	"context"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/config"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/services"
)

// Daemon сервер gRPC с сервисами аутентификации и хранения пользовательских данных
type Daemon struct {
	AuthService   *services.AuthService
	SecretService *services.SecretService
	Address       string
}

// NewDaemon создает новый Daemon с указанными настройками
func NewDaemon(cfg config.Config) *Daemon {
	authService := services.NewAuthService(cfg)
	secretService := services.NewSecretService(cfg)
	return &Daemon{
		AuthService:   authService,
		SecretService: secretService,
		Address:       cfg.GRPC.Address,
	}
}

// Run функция запуска gRPC сервера с сервисами аутентификации и хранения пользовательских данных
func (daemon *Daemon) Run(ctx context.Context) {
	interceptor := interceptors.NewAuthInterceptor(daemon.AuthService.TokenManager)

	services.NewServer(
		daemon.Address,
		services.WithServices(daemon.AuthService, daemon.SecretService),
		services.WithUnaryInterceptors(interceptor.Unary()),
		services.WithStreamInterceptors(interceptor.Stream()),
	).Run(ctx)
}
