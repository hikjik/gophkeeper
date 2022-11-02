package config

// Config содержит настройки клиента
type Config struct {
	GRPC       GRPCConfig       `mapstructure:"grpc"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

// GRPCConfig настройки GRPC
type GRPCConfig struct {
	Address string `mapstructure:"address"`
}

// EncryptionConfig настройки шифрования
type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}
