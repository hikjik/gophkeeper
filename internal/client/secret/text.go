package secret

import "fmt"

var _ Secret = (*Text)(nil)

// Text текстовые данные
type Text struct {
	Data string
}

// Type возвращает тип хранимой информации
func (t Text) Type() string {
	return secretTypeText
}

// String функция отображения приватной информации
func (t Text) String() string {
	return fmt.Sprintf("TextData: %s", t.Data)
}
