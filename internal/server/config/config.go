package config

import "time"

// Config содержит настройки сервера
type Config struct {
	GRPC GRPCConfig    `mapstructure:"grpc"`
	DB   StorageConfig `mapstructure:"db"`
	Auth AuthConfig    `mapstructure:"auth"`
	Hash HashConfig    `mapstructure:"hasher"`
}

// GRPCConfig настройки GRPC
type GRPCConfig struct {
	Address string `mapstructure:"address"`
}

// StorageConfig настройки базы данных сервера
type StorageConfig struct {
	URL string `mapstructure:"url"`
}

// AuthConfig настройки аутентификации
type AuthConfig struct {
	Key            string        `mapstructure:"key"`
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
}

// HashConfig настройки хеширования
type HashConfig struct {
	Key string `mapstructure:"key"`
}
