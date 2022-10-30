package hasher

// Hasher общий интерфейс, реализуемый хэш-функциями
type Hasher interface {
	// Hash вычисляет хэш переданной строки
	Hash(data string) (string, error)
	// IsValid проверяет хэш переданной строки
	IsValid(data string, hash string) (bool, error)
}
