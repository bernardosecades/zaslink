package service

import (
	"github.com/bernardosecades/sharesecret/repository"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetContentSecret(t *testing.T) {
	key := "11111111111111111111111111111111"
	pass := "@myPassword"

	sut := NewSecretService(repository.NewMockSecretRepository(), key, "")
	content, _ := sut.GetContentSecret("10", pass)

	assert.Equal(t, "My name is Bernie", content)
}

func TestCreateSecret(t *testing.T) {
	key := "11111111111111111111111111111111"
	pass := ""
	defaultPwd := "@myPassword"

	sut := NewSecretService(repository.NewMockSecretRepository(), key, defaultPwd)
	secret, _ := sut.CreateSecret("My secret content", pass)

	assert.Equal(t, "My name is Bernie", sut.DecryptContentSecret(secret.Content, defaultPwd))
}
