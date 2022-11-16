package pg

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
)

func newUserMock() (storage.UserStorage, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create sql mock db")
	}
	return &userStorage{db: db}, mock
}

func newTestUser() (*models.User, int) {
	return &models.User{
		Email:        "test@mail.ru",
		PasswordHash: "password_hash",
	}, 0
}

func TestPostgresStorage_GetUser(t *testing.T) {
	s, mock := newUserMock()
	user, userID := newTestUser()

	t.Run("SuccessfulGetUser", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM users WHERE").
			WithArgs(user.Email, user.PasswordHash).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		id, err := s.GetUser(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, id, userID)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM users WHERE").
			WithArgs(user.Email, user.PasswordHash).
			WillReturnError(sql.ErrNoRows)

		_, err := s.GetUser(context.Background(), user)
		assert.ErrorIs(t, err, storage.ErrUserNotFound)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresStorage_PutUser(t *testing.T) {
	s, mock := newUserMock()
	user, userID := newTestUser()

	t.Run("SuccessfulPutUser", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Email, user.PasswordHash).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		id, err := s.PutUser(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, id, userID)
	})
	t.Run("UserExists", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Email, user.PasswordHash).
			WillReturnError(sql.ErrNoRows)

		_, err := s.PutUser(context.Background(), user)
		assert.ErrorIs(t, err, storage.ErrUserConflict)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
