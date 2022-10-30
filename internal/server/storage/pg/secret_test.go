package pg

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
)

func newSecretMock() (storage.SecretStorage, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create sql mock db")
	}
	return &secretStorage{db: db}, mock
}

func TestPostgresStorage_GetSecret(t *testing.T) {
	s, mock := newSecretMock()
	secretName := "TestName"
	secretOwnerID := 0

	t.Run("Secret not exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT content, version FROM secrets WHERE").
			WithArgs(secretName, secretOwnerID).
			WillReturnError(storage.ErrSecretNotFound)

		_, err := s.GetSecret(context.Background(), secretName, secretOwnerID)
		assert.Error(t, err)
		assert.ErrorIs(t, storage.ErrSecretNotFound, err)
	})

	t.Run("Secret exists", func(t *testing.T) {
		versionExpected := uuid.New()
		contentExpected := []byte("TestContent")

		mock.ExpectQuery("SELECT content, version FROM secrets WHERE").
			WithArgs(secretName, secretOwnerID).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"content", "version"}).
					AddRow(contentExpected, versionExpected))

		secret, err := s.GetSecret(context.Background(), secretName, secretOwnerID)
		assert.NoError(t, err)
		assert.Equal(t, versionExpected, secret.Version)
		assert.Equal(t, contentExpected, secret.Content)
	})
}

func TestPostgresStorage_PutSecret(t *testing.T) {
	s, mock := newSecretMock()

	t.Run("Create secret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "TestName",
			Content: []byte("TestContent"),
			Version: uuid.UUID{},
			OwnerID: 0,
		}

		t.Run("Name conflict", func(t *testing.T) {
			mock.ExpectQuery("INSERT INTO secrets").
				WithArgs(secret.Name, secret.Content, secret.OwnerID).
				WillReturnError(storage.ErrSecretNameConflict)

			_, err := s.PutSecret(context.Background(), secret)
			assert.Error(t, err)
			assert.ErrorIs(t, storage.ErrSecretNameConflict, err)
		})

		t.Run("Successful creation", func(t *testing.T) {
			versionExpected := uuid.New()

			mock.ExpectQuery("INSERT INTO secrets").
				WithArgs(secret.Name, secret.Content, secret.OwnerID).
				WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(versionExpected))

			versionActual, err := s.PutSecret(context.Background(), secret)
			assert.NoError(t, err)
			assert.Equal(t, versionExpected, versionActual)
		})
	})

	t.Run("Update secret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "TestName",
			Content: []byte("TestContent"),
			Version: uuid.New(),
			OwnerID: 0,
		}

		t.Run("Secret not found", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT version FROM secrets WHERE").
				WithArgs(secret.Name, secret.OwnerID).
				WillReturnError(storage.ErrSecretNotFound)

			_, err := s.PutSecret(context.Background(), secret)
			assert.Error(t, err)
			assert.ErrorIs(t, storage.ErrSecretNotFound, err)
		})

		t.Run("Version conflict", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT version FROM secrets WHERE").
				WithArgs(secret.Name, secret.OwnerID).
				WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(uuid.New()))

			_, err := s.PutSecret(context.Background(), secret)
			assert.Error(t, err)
			assert.ErrorIs(t, storage.ErrSecretVersionConflict, err)
		})

		t.Run("Error on update", func(t *testing.T) {
			updateError := errors.New("some error")

			mock.ExpectBegin()
			mock.ExpectQuery("SELECT version FROM secrets WHERE").
				WithArgs(secret.Name, secret.OwnerID).
				WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(secret.Version))
			mock.ExpectQuery("UPDATE secrets SET version").
				WithArgs(secret.Content, secret.OwnerID, secret.Name, secret.Version).
				WillReturnError(updateError)
			mock.ExpectRollback()

			_, err := s.PutSecret(context.Background(), secret)
			assert.Error(t, err)
			assert.ErrorIs(t, err, updateError)
		})

		t.Run("Successful update", func(t *testing.T) {
			newVersion := uuid.New()

			mock.ExpectBegin()
			mock.ExpectQuery("SELECT version FROM secrets WHERE").
				WithArgs(secret.Name, secret.OwnerID).
				WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(secret.Version))
			mock.ExpectQuery("UPDATE secrets SET version").
				WithArgs(secret.Content, secret.OwnerID, secret.Name, secret.Version).
				WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(newVersion))
			mock.ExpectCommit()

			newVersionActual, err := s.PutSecret(context.Background(), secret)
			assert.NoError(t, err)
			assert.Equal(t, newVersion, newVersionActual)
		})
	})
}

func TestPostgresStorage_ListSecrets(t *testing.T) {
	s, mock := newSecretMock()
	userID := 0

	t.Run("Select Error", func(t *testing.T) {
		errExpected := errors.New("some error")
		mock.ExpectQuery("SELECT name, content, version FROM secrets").
			WithArgs(userID).WillReturnError(errExpected)

		_, err := s.ListSecrets(context.Background(), userID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errExpected)
	})

	t.Run("Successful list", func(t *testing.T) {
		secrets := []*models.Secret{
			{
				Name:    "Name1",
				Content: []byte("Content1"),
				Version: uuid.New(),
				OwnerID: userID,
			},
			{
				Name:    "Name2",
				Content: []byte("Content2"),
				Version: uuid.New(),
				OwnerID: userID,
			},
		}

		mock.ExpectQuery("SELECT name, content, version FROM secrets").
			WithArgs(userID).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"name", "content", "version"}).
					AddRow(secrets[0].Name, secrets[0].Content, secrets[0].Version).
					AddRow(secrets[1].Name, secrets[1].Content, secrets[1].Version))

		secretsActual, err := s.ListSecrets(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, len(secrets), len(secretsActual))
		for i := 0; i < len(secrets); i++ {
			assert.Equal(t, secrets[i].Name, secretsActual[i].Name)
			assert.Equal(t, secrets[i].Content, secretsActual[i].Content)
			assert.Equal(t, secrets[i].Version, secretsActual[i].Version)
			assert.Equal(t, secrets[i].OwnerID, secretsActual[i].OwnerID)
		}
	})
}
