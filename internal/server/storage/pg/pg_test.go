package pg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	databaseURL := os.Getenv("DB_URL")

	assert.NoError(t, migrate(databaseURL))
	assert.Error(t, migrate(""))
}

func TestNewSecretStorage(t *testing.T) {
	databaseURL := os.Getenv("DB_URL")

	storage, err := NewSecretStorage(databaseURL)
	assert.NotNil(t, storage)
	assert.NoError(t, err)

	_, err = NewSecretStorage("")
	assert.Error(t, err)
}

func TestNewUserStorage(t *testing.T) {
	databaseURL := os.Getenv("DB_URL")

	storage, err := NewUserStorage(databaseURL)
	assert.NotNil(t, storage)
	assert.NoError(t, err)

	_, err = NewUserStorage("")
	assert.Error(t, err)
}
