package models

import "fmt"

var _ Secret = (*Credentials)(nil)

// Credentials пара логин пароль
type Credentials struct {
	Login    string
	Password string
}

// Type возвращает тип хранимой информации
func (c Credentials) Type() SecretType {
	return secretTypeCredentials
}

// String функция отображения приватной информации
func (c Credentials) String() string {
	return fmt.Sprintf("Login: %s, Password: %s", c.Login, c.Password)
}
