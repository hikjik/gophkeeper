package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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

// PutSecret создает или обновляет содержимое секрета secret в базе данных
func (s *secretStorage) PutSecret(ctx context.Context, secret *models.Secret) (uuid.UUID, error) {
	if secret.Version == uuid.Nil {
		row := s.db.QueryRowContext(
			ctx,
			`INSERT INTO secrets (name, content, owner_id)
                   VALUES($1, $2, $3)
                   ON CONFLICT DO NOTHING RETURNING version`,
			secret.Name, secret.Content, secret.OwnerID,
		)
		var newVersion uuid.UUID
		err := row.Scan(&newVersion)
		if errors.Is(err, sql.ErrNoRows) {
			return newVersion, storage.ErrSecretNameConflict
		}
		return newVersion, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}

	row := tx.QueryRowContext(
		ctx,
		`SELECT version FROM secrets WHERE name = ($1) AND owner_id = ($2)`,
		secret.Name, secret.OwnerID,
	)
	var oldVersion uuid.UUID
	err = row.Scan(&oldVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, storage.ErrSecretNotFound
		}
		return uuid.Nil, err
	}
	if oldVersion != secret.Version {
		return uuid.Nil, storage.ErrSecretVersionConflict
	}

	SQLQuery := `
        UPDATE secrets
        SET version = uuid_generate_v4(), content = ($1)
        WHERE owner_id = ($2) AND name = ($3) AND version = ($4)
        RETURNING version`

	row = tx.QueryRowContext(ctx, SQLQuery, secret.Content, secret.OwnerID, secret.Name, secret.Version)
	var newVersion uuid.UUID
	err = row.Scan(&newVersion)
	if err != nil {
		if tx.Rollback() != nil {
			log.Error().Msg("Failed to rollback")
		}
		return uuid.Nil, err
	}

	return newVersion, tx.Commit()
}

// ListSecrets возвращает список всех секретов пользователя с указанным идентификатором
func (s *secretStorage) ListSecrets(ctx context.Context, userID int) ([]*models.Secret, error) {
	rows, err := s.db.QueryContext(
		ctx, `SELECT name, content, version FROM secrets WHERE owner_id = ($1)`, userID)
	if err != nil {
		return nil, err
	}

	secrets := make([]*models.Secret, 0)
	for rows.Next() {
		secret := &models.Secret{
			OwnerID: userID,
		}
		if err = rows.Scan(&secret.Name, &secret.Content, &secret.Version); err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}
