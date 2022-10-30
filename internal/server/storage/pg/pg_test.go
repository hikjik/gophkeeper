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

func TestNew(t *testing.T) {
	databaseURL := os.Getenv("DB_URL")

	storage, err := New(databaseURL)
	assert.NotNil(t, storage)
	assert.NoError(t, err)

	_, err = New("")
	assert.Error(t, err)
}
