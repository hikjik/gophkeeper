package storage

import (
	"context"
	"errors"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
)

// Возможные ошибки при работе с хранилищем UserStorage
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserConflict = errors.New("user conflict")
)

// UserStorage определяет интерфейс для хранения учетных данных пользователей
type UserStorage interface {
	// PutUser сохраняет учетные данные пользователя
	PutUser(ctx context.Context, user *models.User) (*models.User, error)
	// GetUser возвращает ID пользователя с указанными учетными данными
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
}

// Возможные ошибки при работе с хранилищем SecretStorage
var (
	ErrSecretNotFound = errors.New("secret not found")
	ErrSecretConflict = errors.New("secret conflict")
)

// SecretStorage определяет интерфейс для хранения приватных данных пользователей
type SecretStorage interface {
	// GetSecret возвращает секрет с указанным именем name для пользователя c идентификатором userID
	GetSecret(ctx context.Context, name string, userID int) (*models.Secret, error)
	// CreateSecret создает новый секрет
	CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	// UpdateSecret обновляет содержимое секрета
	UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	// DeleteSecret удаляет секрет
	DeleteSecret(ctx context.Context, secret *models.Secret) error
	// ListSecrets возвращает список всех секретов пользователя с указанным идентификатором
	ListSecrets(ctx context.Context, userID int) ([]*models.Secret, error)
}
