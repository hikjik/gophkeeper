package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/token"
)

// TokenGenerator JWT реализация интерфейса token.Generator
type TokenGenerator struct {
	key            []byte
	expirationTime time.Duration
}

var _ token.Generator = (*TokenGenerator)(nil)

// NewJWTTokenGenerator возвращает новый JWT TokenGenerator
func NewJWTTokenGenerator(key string, expirationTime time.Duration) (*TokenGenerator, error) {
	if len(key) < token.MinKeySize {
		return nil, token.ErrInvalidKeySize
	}
	return &TokenGenerator{
		key:            []byte(key),
		expirationTime: expirationTime,
	}, nil
}

// Create возвращает новый JWT токен
func (g *TokenGenerator) Create(userID int) (string, error) {
	payload, err := token.NewPayload(userID, g.expirationTime)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(g.key)
}

// Validate проверяет JWT токен на валидность
func (g *TokenGenerator) Validate(accessToken string) (*token.Payload, error) {
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, token.ErrInvalidToken
		}
		return g.key, nil
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
