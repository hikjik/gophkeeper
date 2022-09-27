package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func contains(array []rune, target rune) bool {
	for _, value := range array {
		if target == value {
			return true
		}
	}
	return false
}

func TestRandomString(t *testing.T) {
	length := 32

	str, err := RandomString(32)
	require.NoError(t, err)
	require.Equal(t, length, len(str))
	for _, r := range str {
		require.True(t, contains(letters, r))
	}
}

func TestRandomString_InvalidLength(t *testing.T) {
	for _, length := range []int{0, -1} {
		_, err := RandomString(length)
		require.ErrorIs(t, err, errInvalidLength)
	}
}
