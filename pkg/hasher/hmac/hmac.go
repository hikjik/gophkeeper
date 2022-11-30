package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/go-developer-ya-practicum/gophkeeper/pkg/hasher"
)

var _ hasher.Hasher = (*Hasher)(nil)

// MinKeySize минимальная длина ключа
const MinKeySize = 16

var (
	// ErrInvalidKeySize ошибка, возникающая при создании Hasher с коротким ключом
	ErrInvalidKeySize = errors.New("invalid key size")
)

// Hasher hmac реализация интерфейса hasher.Hasher
type Hasher struct {
	key []byte
}

// New возвращает новый hmac Hasher
func New(key string) (*Hasher, error) {
	if len(key) < MinKeySize {
		return nil, ErrInvalidKeySize
	}
	return &Hasher{
		key: []byte(key),
	}, nil
}

// Hash вычисляет хэш переданной строки
func (h *Hasher) Hash(data string) (string, error) {
	mac := hmac.New(sha256.New, h.key)
	if _, err := mac.Write([]byte(data)); err != nil {
		return "", err
	}
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum), nil
}

// IsValid проверяет хэш переданной строки
func (h *Hasher) IsValid(data, hash string) (bool, error) {
	mac := hmac.New(sha256.New, h.key)
	if _, err := mac.Write([]byte(data)); err != nil {
		return false, err
	}
	expected, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	actual := mac.Sum(nil)
	return hmac.Equal(expected, actual), nil
}
