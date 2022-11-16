package models

import (
	"encoding/json"
	"errors"
)

// SecretType тип секрета
type SecretType string

const (
	secretTypeCredentials SecretType = "credentials"
	secretTypeText        SecretType = "text"
	secretTypeBin         SecretType = "bin"
	secretTypeCard        SecretType = "card"
)

// Secret приватные данные пользователя
type Secret interface {
	// Type возвращает тип хранимой информации
	Type() SecretType
	// String функция отображения приватной информации
	String() string
}

type container struct {
	Type SecretType      `json:"type"`
	Data json.RawMessage `json:"data"`
}

// DecodeSecret функция декодирования данных пользователя
func DecodeSecret(data []byte) (Secret, error) {
	var c container
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	switch c.Type {
	case secretTypeCredentials:
		var credentials Credentials
		if err := json.Unmarshal(c.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case secretTypeText:
		var text Text
		if err := json.Unmarshal(c.Data, &text); err != nil {
			return nil, err
		}
		return text, nil
	case secretTypeBin:
		var bin Bin
		if err := json.Unmarshal(c.Data, &bin); err != nil {
			return nil, err
		}
		return bin, nil
	case secretTypeCard:
		var card Card
		if err := json.Unmarshal(c.Data, &card); err != nil {
			return nil, err
		}
		return card, nil
	default:
		return nil, errors.New("unknown secret type")
	}
}

// EncodeSecret функция декодирования данных пользователя
func EncodeSecret(secret Secret) ([]byte, error) {
	data, err := json.Marshal(secret)
	if err != nil {
		return nil, err
	}
	return json.Marshal(container{
		Type: secret.Type(),
		Data: data,
	})
}
