package models

// User содержит учетные данные пользователя
type User struct {
	ID           int
	Email        string
	PasswordHash string
}
