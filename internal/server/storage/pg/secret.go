package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
)

type secretStorage struct {
	db *sql.DB
}

var _ storage.SecretStorage = (*secretStorage)(nil)

// NewSecretStorage возвращает объект, реализующий интерфейс storage.SecretStorage
func NewSecretStorage(databaseURL string) (storage.SecretStorage, error) {
	if err := migrate(databaseURL); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	return &secretStorage{db: db}, nil
}

// GetSecret возвращает секрет с указанным именем name для пользователя c идентификатором userID
func (s *secretStorage) GetSecret(ctx context.Context, name string, userID int) (*models.Secret, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT content, version FROM secrets WHERE name = ($1) AND owner_id = ($2)`,
		name, userID,
	)
	secret := &models.Secret{
		Name:    name,
		OwnerID: userID,
	}
	err := row.Scan(&secret.Content, &secret.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrSecretNotFound
	}
	return secret, err
}

// CreateSecret создает новый секрет в базе данных
func (s *secretStorage) CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	row := s.db.QueryRowContext(
		ctx,
		`INSERT INTO secrets (name, content, owner_id)
                   VALUES($1, $2, $3)
                   ON CONFLICT DO NOTHING RETURNING version`,
		secret.Name, secret.Content, secret.OwnerID,
	)
	err := row.Scan(&secret.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return secret, storage.ErrSecretConflict
	}
	return secret, err
}

// UpdateSecret функция обновления секрета в базе данных
func (s *secretStorage) UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	SQLQuery := `
        UPDATE secrets
        SET version = uuid_generate_v4(), content = ($1)
        WHERE owner_id = ($2) AND name = ($3)
        RETURNING version`

	row := s.db.QueryRowContext(ctx, SQLQuery, secret.Content, secret.OwnerID, secret.Name)
	err := row.Scan(&secret.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrSecretNotFound
		}
		return nil, err
	}

	return secret, nil
}

func (s *secretStorage) DeleteSecret(ctx context.Context, secret *models.Secret) error {
	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM secrets WHERE name = ($1) AND owner_id = ($2)`,
		secret.Name,
		secret.OwnerID,
	)
	return err
}

// ListSecrets возвращает список всех секретов пользователя с указанным идентификатором
func (s *secretStorage) ListSecrets(ctx context.Context, userID int) ([]*models.Secret, error) {
	rows, err := s.db.QueryContext(
		ctx, `SELECT name, version, content FROM secrets WHERE owner_id = ($1)`, userID)
	if err != nil {
		return nil, err
	}

	secrets := make([]*models.Secret, 0)
	for rows.Next() {
		secret := &models.Secret{
			OwnerID: userID,
		}
		if err = rows.Scan(&secret.Name, &secret.Version, &secret.Content); err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}
