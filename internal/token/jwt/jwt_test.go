package jwt

import (
	"math/rand"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/token"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/utils"
)

func TestNewJWTTokenGenerator(t *testing.T) {
	for _, length := range []int{8, 16} {
		key, err := utils.RandomString(length)
		require.NoError(t, err)

		generator, err := NewJWTTokenGenerator(key, time.Minute)

		if length >= token.MinKeySize {
			require.NoError(t, err)
			require.NotNil(t, generator)
		} else {
			require.ErrorIs(t, err, token.ErrInvalidKeySize)
		}
	}
}

func TestTokenGenerator_Validate(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		key, err := utils.RandomString(16)
		require.NoError(t, err)

		duration := time.Minute
		generator, err := NewJWTTokenGenerator(key, duration)
		require.NoError(t, err)

		userID := rand.Int()
		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		accessToken, err := generator.Create(userID)
		require.NoError(t, err)
		require.NotEmpty(t, accessToken)

		payload, err := generator.Validate(accessToken)
		require.NoError(t, err)
		require.NotNil(t, payload)
		require.NotZero(t, payload.Id)
		require.Equal(t, payload.UserID, userID)
		require.WithinDuration(t, issuedAt, time.Unix(payload.IssuedAt, 0), time.Second)
		require.WithinDuration(t, expiredAt, time.Unix(payload.ExpiresAt, 0), time.Second)
	})

	t.Run("Expired Token", func(t *testing.T) {
		key, err := utils.RandomString(16)
		require.NoError(t, err)

		generator, err := NewJWTTokenGenerator(key, -time.Minute)
		require.NoError(t, err)

		userID := rand.Int()
		accessToken, err := generator.Create(userID)
		require.NoError(t, err)
		require.NotEmpty(t, accessToken)

		_, err = generator.Validate(accessToken)
		require.Error(t, err)
		require.ErrorIs(t, err, token.ErrExpiredToken)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		key, err := utils.RandomString(16)
		require.NoError(t, err)

		duration := time.Minute
		generator, err := NewJWTTokenGenerator(key, duration)
		require.NoError(t, err)

		userID := rand.Int()
		payload, err := token.NewPayload(userID, duration)
		require.NoError(t, err)

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
		accessToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		_, err = generator.Validate(accessToken)
		require.Error(t, err)
		require.ErrorIs(t, err, token.ErrInvalidToken)
	})
}
