package storage

import (
	"context"
	"errors"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/models"
)

// Возможные ошибки при работе с хранилищем
var (
	ErrEmailIsAlreadyInUse = errors.New("email is already in use")
	ErrInvalidCredentials  = errors.New("invalid request credentials")
)

// Storage определяет интерфейс для хранения приватных данных пользователей
type Storage interface {
	// PutUser сохраняет учетные данные пользователя
	PutUser(ctx context.Context, user *models.User) (userID int, err error)
	// GetUser возвращает ID пользователя с указанными учетными данными
	GetUser(ctx context.Context, user *models.User) (userID int, err error)
}
