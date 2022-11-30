package interceptors

type key int

const (
	// ContextKeyUserID ключ для добавления UserID в контекст при аутентификации
	ContextKeyUserID key = iota
)
