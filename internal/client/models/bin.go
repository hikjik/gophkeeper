package models

var _ Secret = (*Bin)(nil)

// Bin произвольные бинарные данные
type Bin struct {
	Data []byte
}

// Type возвращает тип хранимой информации
func (b Bin) Type() SecretType {
	return secretTypeBin
}

// String функция отображения приватной информации
func (b Bin) String() string {
	return "BINARY DATA"
}
