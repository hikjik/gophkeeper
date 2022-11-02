package generate

import "crypto/rand"

// RandomBytes генерирует последовательность случайных байтов заданной длины
func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
