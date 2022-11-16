package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBin_Type(t *testing.T) {
	assert.Equal(t, secretTypeBin, Bin{}.Type())
}

func TestBin_String(t *testing.T) {
	secret := Bin{Data: []byte("data")}
	assert.Equal(t, "BINARY DATA", secret.String())
}
