package services

import (
	"github.com/bernardosecades/sharesecret/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContentSecret(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"

	sut := NewSecretService(mocks.NewMockSecretRepository(), key)
	content := sut.GetContentSecret("10", pass)

	assert.Equal(t, "My name is Bernie", content)
}

func TestCreateContentSecret(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"

	sut := NewSecretService(mocks.NewMockSecretRepository(), key)
	secret := sut.CreateSecret("My secret content", pass)

	assert.Equal(t, "My secret content", sut.DecryptContentSecret(secret.Content, pass))
}
