package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeSecret(t *testing.T) {
	t.Run("EncodeCredentials", func(t *testing.T) {
		secret := Credentials{
			Login:    "Login",
			Password: "Password",
		}
		expected := []byte(`{"type":"credentials","data":{"Login":"Login","Password":"Password"}}`)

		data, err := EncodeSecret(secret)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	})
	t.Run("EncodeText", func(t *testing.T) {
		secret := Text{Data: "Data"}
		expected := []byte(`{"type":"text","data":{"Data":"Data"}}`)

		data, err := EncodeSecret(secret)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	})
	t.Run("EncodeBin", func(t *testing.T) {
		secret := Bin{Data: []byte("Data")}
		expected := []byte(`{"type":"bin","data":{"Data":"RGF0YQ=="}}`)

		data, err := EncodeSecret(secret)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	})
	t.Run("EncodeCard", func(t *testing.T) {
		secret := Card{
			Number:       "Number",
			ExpiryDate:   "ExpiryDate",
			SecurityCode: "SecurityCode",
			Holder:       "Holder",
		}
		expected := []byte(`{"type":"card","data":{"Number":"Number","ExpiryDate":"ExpiryDate","SecurityCode":"SecurityCode","Holder":"Holder"}}`)

		data, err := EncodeSecret(secret)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	})
}

func TestDecodeSecret(t *testing.T) {
	t.Run("DecodeCredentials", func(t *testing.T) {
		data := []byte(`{"type":"credentials","data":{"Login":"Login","Password":"Password"}}`)

		secret, err := DecodeSecret(data)
		assert.NoError(t, err)
		assert.Equal(t, secret.Type(), secretTypeCredentials)

		credentials, ok := secret.(Credentials)
		assert.True(t, ok)
		assert.Equal(t, "Login", credentials.Login)
		assert.Equal(t, "Password", credentials.Password)
	})
	t.Run("DecodeText", func(t *testing.T) {
		data := []byte(`{"type":"text","data":{"Data":"Data"}}`)

		secret, err := DecodeSecret(data)
		assert.NoError(t, err)
		assert.Equal(t, secret.Type(), secretTypeText)

		text, ok := secret.(Text)
		assert.True(t, ok)
		assert.Equal(t, "Data", text.Data)
	})
	t.Run("DecodeBin", func(t *testing.T) {
		data := []byte(`{"type":"bin","data":{"Data":"RGF0YQ=="}}`)

		secret, err := DecodeSecret(data)
		assert.NoError(t, err)
		assert.Equal(t, secret.Type(), secretTypeBin)

		bin, ok := secret.(Bin)
		assert.True(t, ok)
		assert.Equal(t, []byte("Data"), bin.Data)
	})
	t.Run("DecodeCard", func(t *testing.T) {
		data := []byte(`{"type":"card","data":{"Number":"Number","ExpiryDate":"ExpiryDate","SecurityCode":"SecurityCode","Holder":"Holder"}}`)

		secret, err := DecodeSecret(data)
		assert.NoError(t, err)
		assert.Equal(t, secret.Type(), secretTypeCard)

		card, ok := secret.(Card)
		assert.True(t, ok)
		assert.Equal(t, "Number", card.Number)
		assert.Equal(t, "ExpiryDate", card.ExpiryDate)
		assert.Equal(t, "SecurityCode", card.SecurityCode)
		assert.Equal(t, "Holder", card.Holder)
	})
}
