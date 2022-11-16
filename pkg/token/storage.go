package token

// Storage интерфейс сохранения из загрузки токена
type Storage interface {
	// Load загружает токен
	Load() (accessToken string, err error)
	// Save сохраняет токен
	Save(accessToken string) error
}
