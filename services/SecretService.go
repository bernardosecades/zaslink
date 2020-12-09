package services

import (
	"encoding/hex"
	"github.com/bernardosecades/sharesecret/crypto"
	"github.com/bernardosecades/sharesecret/models"
	"github.com/bernardosecades/sharesecret/repositories"
	uuid "github.com/satori/go.uuid"
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

func (s *SecretService) GetSecret(id string) models.Secret {
	secret, err := s.repository.GetSecret(id)
	if err != nil {
		panic("errornto try to get secret") // TODO handle error better
	}

	return secret
}

func (s *SecretService) GetContentSecret(id string, password string) string {
	secret := s.GetSecret(id)
	return s.DecryptContentSecret(secret.Content, password)
}

func (s *SecretService) CreateSecret(rawContent string, password string) models.Secret {

	if len(password) > 32 {
		panic("password too long")
	}

	return models.Secret{
		uuid.NewV4().String(),
		s.EncryptContentSecret(rawContent, password),
	}
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
