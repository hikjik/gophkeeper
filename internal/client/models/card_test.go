package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCard_Type(t *testing.T) {
	assert.Equal(t, secretTypeCard, Card{}.Type())
}

func TestCard_String(t *testing.T) {
	secret := Card{
		Number:       "Number",
		ExpiryDate:   "ExpiryDate",
		SecurityCode: "SecurityCode",
		Holder:       "Holder",
	}

	expected := "Number: Number, ExpiryDate: ExpiryDate, SecurityCode: SecurityCode, Holder: Holder"
	assert.Equal(t, expected, secret.String())
}
