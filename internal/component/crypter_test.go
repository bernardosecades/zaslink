package component_test

import (
	"github.com/bernardosecades/sharesecret/internal/component"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestEncryptDecryptWithSameKey(t *testing.T) {
	text := []byte("My name is Bernie")
	key := []byte("11111111111111111111111111111111")

	password := "@myPassword"
	copy(key[:], password)

	re, err1 := component.Encrypt(key, text)
	rd, err2 := component.Decrypt(key, re)

	assert.Equal(t, text, rd)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestEncryptDecryptWithDifferentKey(t *testing.T) {
	text := []byte("My name is Bernie")
	key1 := []byte("11111111111111111111111111111111")
	key2 := []byte("11111111111111111111111111111112")

	re, _ := component.Encrypt(key1, text)
	rd, err := component.Decrypt(key2, re)

	assert.NotEqual(t, text, rd)
	assert.Nil(t, rd)
	assert.Equal(t, "cipher: message authentication failed", err.Error())
}
