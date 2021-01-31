package service

import (
	"encoding/hex"

	"github.com/bernardosecades/sharesecret/repository"
	"github.com/bernardosecades/sharesecret/types"

	"errors"
	"time"
)

type SecretService struct {
	repository repository.SecretRepository
	key        string
	defaultPwd string
}

func NewSecretService(r repository.SecretRepository, key string, defaultPwd string) *SecretService {

	if len(key) != 32 {
		panic("key secret should have 32 bytes")
	}

	return &SecretService{r, key, defaultPwd}
}

func (s *SecretService) GetSecret(id string) (types.Secret, error) {
	secret, err := s.repository.GetSecret(id)
	if err != nil {
		return types.Secret{}, err
	}

	err = s.repository.RemoveSecret(id)

	if err != nil {
		return types.Secret{}, err
	}

	return secret, nil
}

func (s *SecretService) HasSecretWithCustomPwd(id string) (bool, error) {

	return s.repository.HasSecretWithCustomPwd(id)
}

func (s *SecretService) GetContentSecret(id string, password string) (string, error) {

	if len(password) == 0 {
		password = s.defaultPwd
	}

	secret, err := s.GetSecret(id)
	if err != nil {
		return "", err
	}

	return s.DecryptContentSecret(secret.Content, password), nil
}

func (s *SecretService) CreateSecret(rawContent string, password string) (types.Secret, error) {

	if len(password) > 32 {
		return types.Secret{}, errors.New("password too long")
	}

	customPwd := true
	if len(password) == 0 {
		customPwd = false
		password = s.defaultPwd
	}

	content := s.EncryptContentSecret(rawContent, password)

	expire := time.Now().UTC().AddDate(0, 0, 5)
	secret, err := s.repository.CreateSecret(content, customPwd, expire)
	if err != nil {
		return types.Secret{}, err
	}

	return secret, nil
}

func (s *SecretService) DecryptContentSecret(content string, password string) string {
	decodeContent, _ := hex.DecodeString(content)
	key := []byte(s.key)
	copy(key[:], password)
	decryptContent, _ := Decrypt(key, decodeContent)

	return string(decryptContent)
}

func (s *SecretService) EncryptContentSecret(content string, password string) string {
	key := []byte(s.key)
	copy(key[:], password)
	encryptContent, _ := Encrypt(key, []byte(content))

	return hex.EncodeToString(encryptContent)
}
