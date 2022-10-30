package token

import (
	"errors"
)

// MinKeySize минимальная длина ключа
const MinKeySize = 16

var (
	// ErrInvalidKeySize генератор использует ключ недостаточной длины
	ErrInvalidKeySize = errors.New("invalid key size")
	// ErrInvalidToken невалидный токен
	ErrInvalidToken = errors.New("token is invalid")
	// ErrExpiredToken закончился срок действия токена
	ErrExpiredToken = errors.New("token has expired")
)

// Manager интерфейс генерации и проверки токенов для аутентификации
type Manager interface {
	// Create создает токен для указанного userID
	Create(userID int) (token string, err error)
	// Validate проверяет токен на валидность
	Validate(accessToken string) (payload *Payload, err error)
}
