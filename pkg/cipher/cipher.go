package cipher

// BlockCipher представляет собой интерфейс с методами шифрования и расшифрования бинарных данных
type BlockCipher interface {
	// Encrypt функция шифрования
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	// Decrypt функция расшифрования
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}
