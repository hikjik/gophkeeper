package server

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/token"
)

// SignUp функция регистрации
func (s *Server) SignUp(ctx context.Context, request *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email or password is empty")
	}

	hash, err := s.Hasher.Hash(request.GetPassword())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to compute hasher")
		return nil, status.Error(codes.Internal, "Failed to compute hasher")
	}
	user := &models.User{
		Email:        request.GetEmail(),
		PasswordHash: hash,
	}

	userID, err := s.UserStorage.PutUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrEmailIsAlreadyInUse) {
			return nil, status.Error(codes.AlreadyExists, "Email is already in use")
		}
		log.Warn().Err(err).Msg("Failed to put user")
		return nil, status.Error(codes.Internal, "Failed to put user")
	}

	accessToken, err := s.TokenManager.Create(userID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate token")
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.SignUpResponse{
		AccessToken: accessToken,
	}, nil
}

// SignIn функция аутентификации
func (s *Server) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email or password is empty")
	}
	hash, err := s.Hasher.Hash(request.GetPassword())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to compute hasher")
		return nil, status.Error(codes.Internal, "Failed to compute hasher")
	}
	user := &models.User{
		Email:        request.GetEmail(),
		PasswordHash: hash,
	}

	userID, err := s.UserStorage.GetUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "Invalid request credentials")
		}
		log.Warn().Err(err).Msg("Failed to get user")
		return nil, status.Error(codes.Internal, "Failed to get user")
	}

	accessToken, err := s.TokenManager.Create(userID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate token")
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.SignInResponse{AccessToken: accessToken}, nil
}

// VerifyToken функция валидации токена
func (s *Server) VerifyToken(ctx context.Context, request *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	accessToken := request.GetAccessToken()
	if accessToken == "" {
		return nil, status.Error(codes.Unauthenticated, "Empty token")
	}

	payload, err := s.TokenManager.Validate(accessToken)
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
