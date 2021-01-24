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
	content, _ := sut.GetContentSecret("10", pass)

	assert.Equal(t, "My name is Bernie", content)
}

func TestCreateSecret(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"

	sut := NewSecretService(mocks.NewMockSecretRepository(), key)
	id, _ := sut.CreateSecret("My secret content", pass)

	assert.Equal(t, "727d7040-aac7-4dc3-ab44-938bfba92ebd", id)
}
