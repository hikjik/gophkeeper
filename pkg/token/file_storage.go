package token

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

// FileStorage файловое хранилище для токена
type FileStorage struct {
	Path string
}

var _ Storage = (*FileStorage)(nil)

// NewFileStorage создает новое файловое хранилище для токена
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{
		Path: path,
	}
}

// Save записывает токен в файл
func (s *FileStorage) Save(accessToken string) error {
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close file with token")
		}
	}()

	_, err = file.WriteString(accessToken)
	return err
}

// Load читает токен из файла
func (s *FileStorage) Load() (string, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		return "", nil
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close file with token")
		}
	}()

	b, err := ioutil.ReadAll(file)
	return string(b), err
}
