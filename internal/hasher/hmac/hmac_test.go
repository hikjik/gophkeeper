package hmac

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/utils"
)

func TestNew(t *testing.T) {
	t.Run("HMAC Hasher", func(t *testing.T) {
		key, err := utils.RandomString(MinKeySize)
		require.NoError(t, err)

		h, err := New(key)
		require.NoError(t, err)
		require.NotNil(t, h)
	})
	t.Run("HMAC Hasher short key", func(t *testing.T) {
		key, err := utils.RandomString(MinKeySize - 1)
		require.NoError(t, err)

		_, err = New(key)
		require.ErrorIs(t, err, ErrInvalidKeySize)
	})
}

func TestHmacHasher_IsValid(t *testing.T) {
	key, err := utils.RandomString(MinKeySize)
	require.NoError(t, err)
	data, err := utils.RandomString(128)
	require.NoError(t, err)

	h, err := New(key)
	require.NoError(t, err)
	require.NotNil(t, h)

	hash, err := h.Hash(data)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	ok, err := h.IsValid(data, hash)
	require.NoError(t, err)
	require.True(t, ok)
}
