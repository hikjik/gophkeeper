package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-developer-ya-practicum/gophkeeper/pkg/token"
)

// TokenManager JWT реализация интерфейса token.Manager
type TokenManager struct {
	key            []byte
	expirationTime time.Duration
}

var _ token.Manager = (*TokenManager)(nil)

// New возвращает новый JWT TokenManager
func New(key string, expirationTime time.Duration) (*TokenManager, error) {
	if len(key) < token.MinKeySize {
		return nil, token.ErrInvalidKeySize
	}
	return &TokenManager{
		key:            []byte(key),
		expirationTime: expirationTime,
	}, nil
}

// Create возвращает новый JWT токен
func (m *TokenManager) Create(userID int) (string, error) {
	payload, err := token.NewPayload(userID, m.expirationTime)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(m.key)
}

// Validate проверяет JWT токен на валидность
func (m *TokenManager) Validate(accessToken string) (*token.Payload, error) {
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, token.ErrInvalidToken
		}
		return m.key, nil
	}

	jwtToken, err := jwt.ParseWithClaims(accessToken, &token.Payload{}, keyFunc)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(validationErr.Inner, token.ErrExpiredToken) {
			return nil, token.ErrExpiredToken
		}
		return nil, token.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*token.Payload)
	if !ok {
		return nil, token.ErrInvalidToken
	}

	return payload, nil
}
