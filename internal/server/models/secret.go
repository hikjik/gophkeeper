package models

import "github.com/google/uuid"

// Secret содержит секретные данные пользователя
type Secret struct {
	Name    string
	Content []byte
	Version uuid.UUID
	OwnerID int
}
