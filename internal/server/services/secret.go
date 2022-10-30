package services

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/config"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage/pg"
)

// SecretService реализация proto.SecretServiceServer
type SecretService struct {
	SecretStorage storage.SecretStorage
	pb.UnimplementedSecretServiceServer
}

var _ pb.SecretServiceServer = (*SecretService)(nil)
var _ Service = (*SecretService)(nil)

// NewSecretService создает новый сервис SecretService
func NewSecretService(cfg config.Config) *SecretService {
	secretStorage, err := pg.NewSecretStorage(cfg.DB.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create storage")
	}

	return &SecretService{SecretStorage: secretStorage}
}

// RegisterService функция регистрации сервиса SecretService на сервере gRPC
func (srv *SecretService) RegisterService(s grpc.ServiceRegistrar) {
	pb.RegisterSecretServiceServer(s, srv)
}

// GetSecret функция получения секрета с указанным в запросе названием
func (srv *SecretService) GetSecret(ctx context.Context, request *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret name is empty")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secret, err := srv.SecretStorage.GetSecret(ctx, request.GetName(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	return &pb.GetSecretResponse{
		Name:    secret.Name,
		Content: secret.Content,
		Version: secret.Version.String(),
	}, nil
}

// SetSecret функция создания или обновления секрета
func (srv *SecretService) SetSecret(ctx context.Context, request *pb.SetSecretRequest) (*pb.SetSecretResponse, error) {
	return nil, nil
}

// ListSecrets возвращает список всех секретов пользователя
func (srv *SecretService) ListSecrets(ctx context.Context, request *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	return nil, nil
}
