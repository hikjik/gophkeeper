package gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"

	"github.com/rs/zerolog/log"

	blockCipher "github.com/go-developer-ya-practicum/gophkeeper/pkg/cipher"
)

const (
	minPasswordSize = 32
	keySize         = 32
	nonceSize       = 12
)

var (
	// ErrInvalidPasswordSize пароль неподходящей длины
	ErrInvalidPasswordSize = errors.New("invalid password size")
)

var _ blockCipher.BlockCipher = (*Cipher)(nil)

// Cipher Блочный шифр AES в режиме GCM
type Cipher struct {
	key   []byte
	nonce []byte
}

// New создает экземпляр Cipher
func New(password string) (*Cipher, error) {
	if len(password) < minPasswordSize {
		return nil, ErrInvalidPasswordSize
	}
	key := sha256.Sum256([]byte(password))
	return &Cipher{
		key:   key[:],
		nonce: key[len(key)-nonceSize:],
	}, nil
}

// Encrypt выполняет шифрование переданной байтовой последовательности
func (c Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(c.key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher")
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher in gcm mode")
		return nil, err
	}

	return aesGCM.Seal(nil, c.nonce, plaintext, nil), nil
}

// Decrypt выполняет расшифрование переданной байтовой последовательности
func (c Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(c.key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher")
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create aes block cipher in gcm mode")
		return nil, err
	}

	return aesGCM.Open(nil, c.nonce, ciphertext, nil)
}
