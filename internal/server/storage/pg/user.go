package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
)

// PutUser сохраняет учетные данные пользователя в базу данных
func (s *postgresStorage) PutUser(ctx context.Context, user *models.User) (int, error) {
	row := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password_hash) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id`,
		user.Email,
		user.PasswordHash,
	)
	var userID int
	err := row.Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, storage.ErrEmailIsAlreadyInUse
	}
	return userID, err
}

// GetUser возвращает ID пользователя с указанными учетными данными
func (s *postgresStorage) GetUser(ctx context.Context, user *models.User) (int, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT id FROM users WHERE email = ($1) AND password_hash = ($2)`,
		user.Email,
		user.PasswordHash,
	)
	var userID int
	err := row.Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, storage.ErrInvalidCredentials
	}
	return userID, err
}
