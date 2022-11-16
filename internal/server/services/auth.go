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
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage/pg"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/hasher"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/hasher/hmac"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/token"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/token/jwt"
)

// AuthService реализация proto.AuthServiceServer
type AuthService struct {
	pb.UnimplementedAuthServiceServer

	UserStorage  storage.UserStorage
	TokenManager token.Manager
	Hasher       hasher.Hasher
}

var _ pb.AuthServiceServer = (*AuthService)(nil)
var _ Service = (*AuthService)(nil)

// NewAuthService создает новый сервис AuthService
func NewAuthService(cfg config.Config) *AuthService {
	userStorage, err := pg.NewUserStorage(cfg.DB.URL)
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

	return &AuthService{
		UserStorage:  userStorage,
		TokenManager: tokenManager,
		Hasher:       hmacHasher,
	}
}

// RegisterService функция регистрации сервиса AuthService на сервере gRPC
func (srv *AuthService) RegisterService(s grpc.ServiceRegistrar) {
	pb.RegisterAuthServiceServer(s, srv)
}

// SignUp функция регистрации
func (srv *AuthService) SignUp(ctx context.Context, request *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email or password is empty")
	}

	hash, err := srv.Hasher.Hash(request.GetPassword())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to compute hasher")
		return nil, status.Error(codes.Internal, "Failed to compute hasher")
	}

	user, err := srv.UserStorage.PutUser(ctx, &models.User{
		Email:        request.GetEmail(),
		PasswordHash: hash,
	})
	if err != nil {
		if errors.Is(err, storage.ErrUserConflict) {
			return nil, status.Error(codes.AlreadyExists, "Email is already in use")
		}
		log.Warn().Err(err).Msg("Failed to put user")
		return nil, status.Error(codes.Internal, "Failed to put user")
	}

	accessToken, err := srv.TokenManager.Create(user.ID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate token")
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.SignUpResponse{
		AccessToken: accessToken,
	}, nil
}

// SignIn функция аутентификации
func (srv *AuthService) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email or password is empty")
	}
	hash, err := srv.Hasher.Hash(request.GetPassword())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to compute hasher")
		return nil, status.Error(codes.Internal, "Failed to compute hasher")
	}

	user, err := srv.UserStorage.GetUser(ctx, &models.User{
		Email:        request.GetEmail(),
		PasswordHash: hash,
	})
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, "Invalid request credentials")
		}
		log.Warn().Err(err).Msg("Failed to get user")
		return nil, status.Error(codes.Internal, "Failed to get user")
	}

	accessToken, err := srv.TokenManager.Create(user.ID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate token")
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.SignInResponse{AccessToken: accessToken}, nil
}

// VerifyToken функция валидации токена
func (srv *AuthService) VerifyToken(ctx context.Context, request *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	accessToken := request.GetAccessToken()
	if accessToken == "" {
		return nil, status.Error(codes.Unauthenticated, "Empty token")
	}

	payload, err := srv.TokenManager.Validate(accessToken)
	if err != nil {
		if errors.Is(err, token.ErrExpiredToken) {
			return nil, status.Error(codes.Unauthenticated, "Token Expired")
		}
		if errors.Is(err, token.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "Token invalid")
		}
		return nil, status.Error(codes.Internal, "Failed to validate token")
	}

	return &pb.VerifyTokenResponse{
		UserId: int32(payload.UserID),
	}, nil
}
