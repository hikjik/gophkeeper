package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText_Type(t *testing.T) {
	assert.Equal(t, secretTypeText, Text{}.Type())
}

func TestText_String(t *testing.T) {
	secret := Text{Data: "data"}
	assert.Equal(t, "TextData: data", secret.String())
}
