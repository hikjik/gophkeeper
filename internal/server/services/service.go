package services

import "google.golang.org/grpc"

// Service интерфейс сервиса gRPC сервера
type Service interface {
	RegisterService(grpc.ServiceRegistrar)
}
