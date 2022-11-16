package models

import (
	"fmt"
)

var _ Secret = (*Card)(nil)

// Card данные банковской карты
type Card struct {
	Number       string
	ExpiryDate   string
	SecurityCode string
	Holder       string
}

// Type возвращает тип хранимой информации
func (c Card) Type() SecretType {
	return secretTypeCard
}

// String функция отображения приватной информации
func (c Card) String() string {
	return fmt.Sprintf("Number: %s, ExpiryDate: %s, SecurityCode: %s, Holder: %s",
		c.Number, c.ExpiryDate, c.SecurityCode, c.Holder)
}
