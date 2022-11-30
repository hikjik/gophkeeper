package services

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Server сервер gRPC
type Server struct {
	Address            string
	Services           []Service
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

// Option определяет настройки gRPC сервера
type Option func(*Server)

// WithUnaryInterceptors возвращает Option, определяющую функции-перехватчики для одиночных RPC запросов
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(server *Server) {
		server.UnaryInterceptors = append(server.UnaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors возвращает Option, определяющую функции-перехватчики для потоковых RPC запросов
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(server *Server) {
		server.StreamInterceptors = append(server.StreamInterceptors, interceptors...)
	}
}

// WithServices возвращает Option, определяющую сервисы gRPC сервера
func WithServices(services ...Service) Option {
	return func(server *Server) {
		server.Services = append(server.Services, services...)
	}
}

// NewServer создает сервер gRPC с заданными настройками
func NewServer(address string, options ...Option) *Server {
	server := &Server{Address: address}

	for _, option := range options {
		option(server)
	}

	return server
}

// Run устанавливает перехватчики, регистрирует сервисы и запускает gRPC сервер
func (s *Server) Run(ctx context.Context) {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start grpc server")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(s.StreamInterceptors...),
	)

	for _, service := range s.Services {
		service.RegisterService(grpcServer)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(listen); err != nil {
		log.Error().Err(err).Msg("Error on grpc server Serve")
	}
}
