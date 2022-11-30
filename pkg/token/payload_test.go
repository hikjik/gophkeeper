package token

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPayload(t *testing.T) {
	userID := rand.Int()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	payload, err := NewPayload(userID, duration)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, payload.UserID, userID)
	require.WithinDuration(t, issuedAt, time.Unix(payload.IssuedAt, 0), time.Second)
	require.WithinDuration(t, expiredAt, time.Unix(payload.ExpiresAt, 0), time.Second)
}

func TestPayload_Valid(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		payload, err := NewPayload(rand.Int(), time.Minute)

		require.NoError(t, err)
		require.NotNil(t, payload)
		require.NoError(t, payload.Valid())
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		payload, err := NewPayload(rand.Int(), -time.Minute)

		require.NoError(t, err)
		require.NotNil(t, payload)
		require.ErrorIs(t, payload.Valid(), ErrExpiredToken)
	})
}
