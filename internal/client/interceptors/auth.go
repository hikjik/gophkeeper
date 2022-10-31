package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor клиентский перехватчик для авторизации
type AuthInterceptor struct {
	token string
}

// NewAuthInterceptor создает новый перехватчик AuthInterceptor
func NewAuthInterceptor(token string) *AuthInterceptor {
	return &AuthInterceptor{token: token}
}

// Unary возвращает клиентскую функцию-перехватчик для авторизации одиночных RPC запросов
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+interceptor.token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
