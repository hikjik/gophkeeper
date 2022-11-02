package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const size = 32

func TestRandomBytes(t *testing.T) {
	b1, err := RandomBytes(size)
	assert.NoError(t, err)
	assert.Equal(t, size, len(b1))

	b2, err := RandomBytes(size)
	assert.NoError(t, err)
	assert.Equal(t, size, len(b2))

	assert.NotEqual(t, b1, b2)
}