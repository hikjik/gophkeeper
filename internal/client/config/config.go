package config

// Config содержит настройки сервера
type Config struct {
	GRPC GRPCConfig `mapstructure:"grpc"`
}

// GRPCConfig настройки GRPC
type GRPCConfig struct {
	Address string `mapstructure:"address"`
}
