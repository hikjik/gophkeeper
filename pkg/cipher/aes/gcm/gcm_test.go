package gcm

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("ShortPassword", func(t *testing.T) {
		_, err := New("short")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPasswordSize)
	})

	t.Run("PasswordLongEnough", func(t *testing.T) {
		cipher, err := New("11112222333344445555666677778888")
		assert.NoError(t, err)
		assert.NotNil(t, cipher)
		assert.Equal(t, keySize, len(cipher.key))
		assert.Equal(t, nonceSize, len(cipher.nonce))
	})
}

// https://boringssl.googlesource.com/boringssl/+/refs/heads/2564/crypto/cipher/test/cipher_test.txt
var tests = []struct {
	key        string
	nonce      string
	plaintext  string
	ciphertext string
	tag        string
}{
	{
		key:        "0000000000000000000000000000000000000000000000000000000000000000",
		nonce:      "000000000000000000000000",
		plaintext:  "00000000000000000000000000000000",
		ciphertext: "cea7403d4d606b6e074ec5d3baf39d18",
		tag:        "d0d1c8a799996bf0265b98b5d48ab919",
	},
	{
		key:        "feffe9928665731c6d6a8f9467308308feffe9928665731c6d6a8f9467308308",
		nonce:      "cafebabefacedbaddecaf888",
		plaintext:  "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b391aafd255",
		ciphertext: "522dc1f099567d07f47f37a32a84427d643a8cdcbfe5c0c97598a2bd2555d1aa8cb08e48590dbb3da7b08b1056828838c5f61e6393ba7a0abcc9f662898015ad",
		tag:        "b094dac5d93471bdec1a502270e3cc6c",
	},
}

func hexDecodeString(t *testing.T, s string) []byte {
	t.Helper()

	b, err := hex.DecodeString(s)
	assert.NoError(t, err)
	return b
}

func TestCipher_Encrypt(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("EncryptTest%d", i), func(t *testing.T) {
			cipher := &Cipher{
				key:   hexDecodeString(t, tt.key),
				nonce: hexDecodeString(t, tt.nonce),
			}

			plaintext := hexDecodeString(t, tt.plaintext)
			ciphertextExpected := hexDecodeString(t, tt.ciphertext+tt.tag)
			ciphertextActual, err := cipher.Encrypt(plaintext)
			assert.NoError(t, err)
			assert.Equal(t, ciphertextExpected, ciphertextActual)
		})
	}
}

func TestCipher_Decrypt(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("DecryptTest%d", i), func(t *testing.T) {
			cipher := &Cipher{
				key:   hexDecodeString(t, tt.key),
				nonce: hexDecodeString(t, tt.nonce),
			}

			plaintextExpected := hexDecodeString(t, tt.plaintext)
			ciphertext := hexDecodeString(t, tt.ciphertext+tt.tag)
			plaintextActual, err := cipher.Decrypt(ciphertext)
			assert.NoError(t, err)
			assert.Equal(t, plaintextExpected, plaintextActual)
		})
	}
}
