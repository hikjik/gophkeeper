package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/config"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
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

// GetSecret возвращает приватные данные пользователя по указанному в запросе названию
func (srv *SecretService) GetSecret(
	ctx context.Context,
	request *pb.GetSecretRequest,
) (*pb.GetSecretResponse, error) {
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

// CreateSecret сохраняет на сервере новые приватные данные пользователя
func (srv *SecretService) CreateSecret(
	ctx context.Context,
	request *pb.CreateSecretRequest,
) (*pb.CreateSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(request.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secret := &models.Secret{
		Name:    request.GetName(),
		Content: request.GetContent(),
		Version: uuid.UUID{},
		OwnerID: userID,
	}
	version, err := srv.SecretStorage.CreateSecret(ctx, secret)
	if err != nil {
		if errors.Is(err, storage.ErrSecretNameConflict) {
			return nil, status.Error(codes.AlreadyExists, "secret already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &pb.CreateSecretResponse{
		Name:    request.GetName(),
		Version: version.String(),
	}, nil
}

// UpdateSecret обновляет приватные данные пользователя
func (srv *SecretService) UpdateSecret(
	ctx context.Context,
	request *pb.UpdateSecretRequest,
) (*pb.UpdateSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(request.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}
	version, err := uuid.Parse(request.GetVersion())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid secret version")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secret := &models.Secret{
		Name:    request.GetName(),
		Content: request.GetContent(),
		Version: version,
		OwnerID: userID,
	}

	newVersion, err := srv.SecretStorage.UpdateSecret(ctx, secret)
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		if errors.Is(err, storage.ErrSecretVersionConflict) {
			return nil, status.Error(codes.NotFound, "secret version conflict")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &pb.UpdateSecretResponse{
		Name:    request.GetName(),
		Version: newVersion.String(),
	}, nil
}

// DeleteSecret удаляет приватные данные пользователя
func (srv *SecretService) DeleteSecret(
	ctx context.Context,
	request *pb.DeleteSecretRequest,
) (*pb.DeleteSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	version, err := uuid.Parse(request.GetVersion())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid secret version")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secret := &models.Secret{
		Name:    request.GetName(),
		Version: version,
		OwnerID: userID,
	}

	if err = srv.SecretStorage.DeleteSecret(ctx, secret); err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		if errors.Is(err, storage.ErrSecretVersionConflict) {
			return nil, status.Error(codes.NotFound, "secret version conflict")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &pb.DeleteSecretResponse{
		Name: request.GetName(),
	}, nil
}

// FetchSecrets возвращает список секретов пользователя, не загружая их контент
func (srv *SecretService) FetchSecrets(
	ctx context.Context,
	_ *pb.FetchSecretsRequest,
) (*pb.FetchSecretsResponse, error) {
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secrets, err := srv.SecretStorage.ListSecrets(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list secrets")
	}

	pbSecrets := make([]*pb.SecretInfo, 0, len(secrets))
	for _, secret := range secrets {
		pbSecrets = append(pbSecrets, &pb.SecretInfo{
			Name:    secret.Name,
			Version: secret.Version.String(),
		})
	}
	return &pb.FetchSecretsResponse{
		Secrets: pbSecrets,
	}, nil
}
