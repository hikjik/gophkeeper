package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var errInvalidLength = errors.New("invalid length")

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString генерирует случайную строку заданной длины из символов letters
func RandomString(length int) (string, error) {
	if length <= 0 {
		return "", errInvalidLength
	}

	var builder strings.Builder
	builder.Grow(length)
	for i := 0; i < length; i++ {
		index := rand.Intn(len(letters))
		builder.WriteRune(letters[index])
	}
	return builder.String(), nil
}
