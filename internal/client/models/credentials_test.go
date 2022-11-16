package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentials_Type(t *testing.T) {
	assert.Equal(t, secretTypeCredentials, Credentials{}.Type())
}

func TestCredentials_String(t *testing.T) {
	secret := Credentials{Login: "Login", Password: "Password"}
	assert.Equal(t, "Login: Login, Password: Password", secret.String())
}
