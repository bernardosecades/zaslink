package services

import (
	"encoding/hex"
	"github.com/bernardosecades/sharesecret/crypto"
	"github.com/bernardosecades/sharesecret/models"
	"github.com/bernardosecades/sharesecret/repositories"
)

type SecretService struct {
	repository repositories.SecretRepository
	key        string
}

func NewSecretService(r repositories.SecretRepository, key string) SecretService {

	if len(key) != 32 {
		panic("key secret should have 32 bytes")
	}

	return SecretService{r, key}
}

func (s *SecretService) GetSecret(id string) (models.Secret, error) {
	secret, err := s.repository.GetSecret(id)
	if err != nil {
		return models.Secret{}, err
	}

	err = s.repository.UpdateToViewed(id)

	if err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (s *SecretService) GetContentSecret(id string, password string) (string, error) {
	secret, err := s.GetSecret(id)
	if err != nil {
		return "", err
	}

	return s.DecryptContentSecret(secret.Content, password), nil
}

func (s *SecretService) CreateSecret(rawContent string, password string) (models.Secret, error) {

	if len(password) > 32 {
		panic("password too long")
	}

	content := s.EncryptContentSecret(rawContent, password)
	secret, err := s.repository.CreateSecret(content)
	if err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (s *SecretService) DecryptContentSecret(content string, password string) string {
	decodeContent, _ := hex.DecodeString(content)
	key := []byte(s.key)
	copy(key[:], password)
	decryptContent, _ := crypto.Decrypt(key, decodeContent)

	return string(decryptContent)
}

func (s *SecretService) EncryptContentSecret(content string, password string) string {
	key := []byte(s.key)
	copy(key[:], password)
	encryptContent, _ := crypto.Encrypt(key, []byte(content))

	return hex.EncodeToString(encryptContent)
}
