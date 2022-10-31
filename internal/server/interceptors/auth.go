package interceptors

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/token"
)

// AuthInterceptor серверный перехватчик для авторизации/аутентификации
type AuthInterceptor struct {
	tokenManager token.Manager
}

// NewAuthInterceptor создает новый перехватчик AuthInterceptor
func NewAuthInterceptor(tokenManager token.Manager) *AuthInterceptor {
	return &AuthInterceptor{tokenManager: tokenManager}
}

// Unary возвращает серверную функцию-перехватчик для аутентификации одиночных RPC запросов
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream возвращает серверную функцию-перехватчик для аутентификации потоковых RPC запросов
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		type serverStream struct {
			grpc.ServerStream
			ctx context.Context
		}
		return handler(srv, &serverStream{
			ServerStream: stream,
			ctx:          ctx,
		})
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	publicMethods := []string{"SignIn", "SignUp", "VerifyToken"}

	for _, publicMethod := range publicMethods {
		if strings.HasSuffix(method, publicMethod) {
			return ctx, nil
		}
	}

	var accessToken string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("authorization")
		if len(values) > 0 {
			accessToken = values[0]
		}
	}

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	if len(accessToken) == 0 {
		return nil, status.Error(codes.Unauthenticated, "empty token")
	}

	payload, err := interceptor.tokenManager.Validate(accessToken)
	if err != nil {
		if errors.Is(err, token.ErrExpiredToken) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}
		if errors.Is(err, token.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "token invalid")
		}
		return nil, status.Error(codes.Internal, "failed to validate token")
	}

	return context.WithValue(ctx, ContextKeyUserID, payload.UserID), nil
}
