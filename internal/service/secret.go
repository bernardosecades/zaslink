package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bernardosecades/sharesecret/internal/entity"
	customEvents "github.com/bernardosecades/sharesecret/internal/events"
	"github.com/bernardosecades/sharesecret/internal/expiration"
	"github.com/bernardosecades/sharesecret/internal/repository"
	"github.com/bernardosecades/sharesecret/pkg/crypter"
	"github.com/bernardosecades/sharesecret/pkg/events"
	"github.com/google/uuid"
)

const (
	keyLength        = 32
	maxPassLength    = 12
	minPassLength    = 4
	maxContentLength = 10000
)

var (
	ErrContentEmpty       = errors.New("content cannot be empty")
	ErrMissingID          = errors.New("missing id")
	ErrContentTooLong     = errors.New("content too long")
	ErrPassTooLong        = errors.New("pass too long")
	ErrPassTooShort       = errors.New("pass too short")
	ErrSecretDoesNotExist = errors.New("secret does not exist")
	ErrInvalidPassword    = errors.New("invalid password")
)

type SecretRepository interface {
	GetSecret(ctx context.Context, id string) (entity.Secret, error)
	SaveSecret(ctx context.Context, secret entity.Secret) error
	DeleteSecret(ctx context.Context, privateID string) (entity.Secret, error)
}

type SecretService struct {
	secretRepo SecretRepository
	publisher  events.Publisher[map[string]string]
	defaultPwd string
	key        string
}

// NewSecretService constructor
func NewSecretService(secretRepo SecretRepository, publisher events.Publisher[map[string]string], defaultPwd, key string) *SecretService {
	if secretRepo == nil {
		panic("secretRepo cannot be nil")
	}
	if publisher == nil {
		panic("publisher cannot be nil")
	}

	if len(key) != keyLength {
		panic("key must be 32 bytes")
	}
	if len(defaultPwd) > maxPassLength {
		panic(fmt.Sprintf("defaultPwd must be <= %d bytes", maxPassLength))
	}

	if len(defaultPwd) < minPassLength {
		panic(fmt.Sprintf("defaultPwd must be >= %d bytes", minPassLength))
	}

	return &SecretService{secretRepo: secretRepo, publisher: publisher, defaultPwd: defaultPwd, key: key}
}

// CreateSecret create handler method
func (s *SecretService) CreateSecret(ctx context.Context, content, pwd, expirationStr string) (entity.Secret, error) {
	if pwd == "" {
		return s.createSecret(ctx, []byte(content), s.defaultPwd, expirationStr, false)
	}
	return s.createSecret(ctx, []byte(content), pwd, expirationStr, true)
}

// DeleteSecret delete handler method
func (s *SecretService) DeleteSecret(ctx context.Context, privateID string) error {
	secret, err := s.secretRepo.DeleteSecret(ctx, privateID)
	if err != nil {
		if errors.Is(err, repository.ErrSecretNotFound) {
			return ErrSecretDoesNotExist
		}
		return err
	}

	go func() {
		// We use empty context instead of ctx because maybe the context was cancelled (Example: client close the connection, request is cancelled in http/2 or the response has been written back to the client)
		_ = s.publisher.Publish(context.Background(), customEvents.NewSecretDeleted(secret))
	}()

	return nil
}

// RetrieveSecret retrieve handler method
func (s *SecretService) RetrieveSecret(ctx context.Context, ID string, pwd string) (entity.Secret, error) {
	if pwd == "" {
		return s.retrieveSecret(ctx, ID, s.defaultPwd)
	}
	return s.retrieveSecret(ctx, ID, pwd)
}

func (s *SecretService) retrieveSecret(ctx context.Context, ID string, pwd string) (entity.Secret, error) {
	if ID == "" {
		return entity.Secret{}, ErrMissingID
	}
	if len(pwd) > maxPassLength {
		return entity.Secret{}, ErrPassTooLong
	}
	if len(pwd) < minPassLength {
		return entity.Secret{}, ErrPassTooShort
	}

	secret, err := s.secretRepo.GetSecret(ctx, ID)
	if err != nil {
		if errors.Is(err, repository.ErrSecretNotFound) {
			return entity.Secret{}, ErrSecretDoesNotExist
		}
		return entity.Secret{}, err
	}

	decryptContent, err := s.decryptContent(secret.Content, pwd)
	if err != nil {
		return entity.Secret{}, ErrInvalidPassword
	}

	secret.Viewed = true
	secret.UpdatedAt = time.Now().UTC()
	err = s.secretRepo.SaveSecret(ctx, secret)
	if err != nil {
		return entity.Secret{}, fmt.Errorf("could not save handler for ID %s: %w", ID, err)
	}

	secret.Content = decryptContent

	go func() {
		// We use empty context instead of ctx because maybe the context was cancelled (Example: client close the connection, request is cancelled in http/2 or the response has been written back to the client)
		_ = s.publisher.Publish(context.Background(), customEvents.NewSecretViewed(secret))
	}()

	return secret, nil
}

func (s *SecretService) createSecret(ctx context.Context, content []byte, pwd, expirationStr string, customPwd bool) (entity.Secret, error) {
	if expirationStr == "" {
		expirationStr = string(expiration.OneDay)
	}

	exp, err := expiration.ValidateExpiration(expirationStr)
	if err != nil {
		return entity.Secret{}, err
	}

	if len(content) == 0 {
		return entity.Secret{}, ErrContentEmpty
	}

	if len(content) > maxContentLength {
		return entity.Secret{}, ErrContentTooLong
	}

	if len(pwd) > maxPassLength {
		return entity.Secret{}, ErrPassTooLong
	}

	if len(pwd) < minPassLength {
		return entity.Secret{}, ErrPassTooShort
	}

	contentEncrypted, err := s.encryptContent(content, pwd)
	if err != nil {
		return entity.Secret{}, err
	}

	now := time.Now().UTC()
	secret := entity.Secret{
		ID:        uuid.New().String(),
		PrivateID: uuid.New().String(),
		Content:   contentEncrypted,
		CustomPwd: customPwd,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiredAt: now.Add(exp.Duration()),
	}

	err = s.secretRepo.SaveSecret(ctx, secret)
	if err != nil {
		return entity.Secret{}, fmt.Errorf("could not save content: %w", err)
	}

	go func() {
		// We use empty context instead of ctx because maybe the context was cancelled (Example: client close the connection, request is cancelled in http/2 or the response has been written back to the client)
		_ = s.publisher.Publish(context.Background(), customEvents.NewSecretCreated(secret))
	}()

	return secret, nil
}

func (s *SecretService) encryptContent(content []byte, pwd string) ([]byte, error) {
	key := crypter.MergePwdIntoKey(s.key, pwd)
	contentEncrypted, err := crypter.Encrypt(key, []byte(content))

	if err != nil {
		return nil, fmt.Errorf("could not encrypt content: %w", err)
	}
	return contentEncrypted, nil
}

func (s *SecretService) decryptContent(content []byte, pwd string) ([]byte, error) {
	key := crypter.MergePwdIntoKey(s.key, pwd)
	decryptContent, err := crypter.Decrypt(key, content)
	if err != nil {
		return nil, fmt.Errorf("could not deecrypt content: %w", err)
	}
	return decryptContent, nil
}
