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

func TestNewJWTTokenManager(t *testing.T) {
	for _, length := range []int{8, 16} {
		key, err := utils.RandomString(length)
		require.NoError(t, err)

		manager, err := New(key, time.Minute)

		if length >= token.MinKeySize {
			require.NoError(t, err)
			require.NotNil(t, manager)
		} else {
			require.ErrorIs(t, err, token.ErrInvalidKeySize)
		}
	}
}

func TestTokenManager_Validate(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		key, err := utils.RandomString(16)
		require.NoError(t, err)

		duration := time.Minute
		manager, err := New(key, duration)
		require.NoError(t, err)

		userID := rand.Int()
		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		accessToken, err := manager.Create(userID)
		require.NoError(t, err)
		require.NotEmpty(t, accessToken)

		payload, err := manager.Validate(accessToken)
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

		manager, err := New(key, -time.Minute)
		require.NoError(t, err)

		userID := rand.Int()
		accessToken, err := manager.Create(userID)
		require.NoError(t, err)
		require.NotEmpty(t, accessToken)

		_, err = manager.Validate(accessToken)
		require.Error(t, err)
		require.ErrorIs(t, err, token.ErrExpiredToken)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		key, err := utils.RandomString(16)
		require.NoError(t, err)

		duration := time.Minute
		manager, err := New(key, duration)
		require.NoError(t, err)

		userID := rand.Int()
		payload, err := token.NewPayload(userID, duration)
		require.NoError(t, err)

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
		accessToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		_, err = manager.Validate(accessToken)
		require.Error(t, err)
		require.ErrorIs(t, err, token.ErrInvalidToken)
	})
}
