//go:build unit

package crypter_test

import (
	"github.com/bernardosecades/zaslink/pkg/crypter"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestEncryptDecryptWithSameKey(t *testing.T) {
	t.Parallel()
	text := []byte("My name is Bernie")
	key := []byte("11111111111111111111111111111111")

	password := "@myPassword"
	copy(key[:], password)

	re, err1 := crypter.Encrypt(key, text)
	rd, err2 := crypter.Decrypt(key, re)

	assert.Equal(t, text, rd)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestEncryptDecryptWithDifferentKey(t *testing.T) {
	t.Parallel()
	text := []byte("My name is Bernie")
	key1 := []byte("11111111111111111111111111111111")
	key2 := []byte("11111111111111111111111111111112")

	re, _ := crypter.Encrypt(key1, text)
	rd, err := crypter.Decrypt(key2, re)

	assert.NotEqual(t, text, rd)
	assert.Nil(t, rd)
	assert.Equal(t, "cipher: message authentication failed", err.Error())
}
