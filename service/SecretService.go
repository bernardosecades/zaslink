package service

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/bernardosecades/sharesecret/repository"
	"github.com/bernardosecades/sharesecret/types"
)

// All errors reported by the service
var (
	ErrSecretNotFound = errors.New("secret not found, it either never existed or has already been viewed")
	ErrNoPassRequired = errors.New("you need a password to see the secret")
	ErrMissingPass    = errors.New("you need a password to see the secret")
	ErrEmptyContent   = errors.New("empty content")
	ErrTextTooLong    = errors.New("text too long")
	ErrPassTooLong    = errors.New("password too long")
)

type SecretService struct {
	repository repository.SecretRepository
	key        string
	defaultPwd string
}

func NewSecretService(r repository.SecretRepository, key string, defaultPwd string) SecretService {

	if len(key) != 32 {
		panic("key secret should have 32 bytes")
	}

	return SecretService{r, key, defaultPwd}
}

func (s *SecretService) GetContentSecret(id string, password string) (string, error) {

	hasPass, err:= s.hasSecretWithCustomPwd(id)

	if err != nil {
		return "", ErrSecretNotFound
	}

	if hasPass && len(password) == 0 {
		return "", ErrMissingPass
	}

	if !hasPass && len(password) > 0 {
		return "", ErrNoPassRequired
	}

	if len(password) == 0 {
		password = s.defaultPwd
	}

	secret, err := s.getSecret(id)
	if err != nil {
		return "", ErrSecretNotFound
	}

	return  s.decryptContentSecret(secret.Content, password)
}

func (s *SecretService) CreateSecret(rawContent string, password string) (types.Secret, error) {

	if len(rawContent) == 0 {
		return types.Secret{}, ErrEmptyContent
	}

	if len(rawContent) > 10000 {
		return types.Secret{}, ErrTextTooLong
	}

	if len(password) > 32 {
		return types.Secret{}, ErrPassTooLong
	}

	customPwd := true
	if len(password) == 0 {
		customPwd = false
		password = s.defaultPwd
	}

	content, err := s.encryptContentSecret(rawContent, password)

	if err != nil {
		return types.Secret{}, err
	}

	expire := time.Now().UTC().AddDate(0, 0, 5)
	secret, err := s.repository.CreateSecret(content, customPwd, expire)
	if err != nil {
		return types.Secret{}, err
	}

	return secret, nil
}

func (s SecretService) hasSecretWithCustomPwd(id string) (bool, error) {

	return s.repository.HasSecretWithCustomPwd(id)
}

func (s SecretService) getSecret(id string) (types.Secret, error) {
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

func (s *SecretService) decryptContentSecret(content string, password string) (string, error) {
	decodeContent, _ := hex.DecodeString(content)
	key := []byte(s.key)
	copy(key[:], password)
	decryptContent, err := Decrypt(key, decodeContent)

	if err != nil {
		return "", err
	}

	return string(decryptContent), nil
}

func (s *SecretService) encryptContentSecret(content string, password string) (string, error) {
	key := []byte(s.key)
	copy(key[:], password)
	encryptContent, err := Encrypt(key, []byte(content))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encryptContent), nil
}
