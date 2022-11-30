package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
)

type userStorage struct {
	db *sql.DB
}

var _ storage.UserStorage = (*userStorage)(nil)

// NewUserStorage возвращает объект, реализующий интерфейс storage.UserStorage
func NewUserStorage(databaseURL string) (storage.UserStorage, error) {
	if err := migrate(databaseURL); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	return &userStorage{db: db}, nil
}

// PutUser сохраняет учетные данные пользователя в базу данных
func (s *userStorage) PutUser(ctx context.Context, user *models.User) (*models.User, error) {
	row := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password_hash) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id`,
		user.Email,
		user.PasswordHash,
	)
	err := row.Scan(&user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrUserConflict
	}
	return user, err
}

// GetUser возвращает ID пользователя с указанными учетными данными
func (s *userStorage) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT id FROM users WHERE email = ($1) AND password_hash = ($2)`,
		user.Email,
		user.PasswordHash,
	)
	err := row.Scan(&user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrUserNotFound
	}
	return user, err
}
