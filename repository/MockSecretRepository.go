package repository

import (
	"github.com/bernardosecades/sharesecret/types"

	"time"
)

type MockSecretRepository struct {
}

func NewMockSecretRepository() *MockSecretRepository {
	return &MockSecretRepository{}
}

// Content: "My secret content" With Key: "11111111111111111111111111111111" And Pass: "@myPassword"
func (s *MockSecretRepository) GetSecret(id string) (types.Secret, error) {
	return types.Secret{
		"727d7040-aac7-4dc3-ab44-938bfba92ebd",
		"cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922",
		false,
		time.Now(),
		time.Now(),
	}, nil
}

func (r *MockSecretRepository) CreateSecret(content string, customPwd bool, expire time.Time) (types.Secret, error) {
	return types.Secret{
		"727d7040-aac7-4dc3-ab44-938bfba92ebd",
		"cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922",
		false,
		time.Now(),
		time.Now(),
	}, nil
}

func (r *MockSecretRepository) RemoveSecret(id string) error {
	return nil
}

func (r *MockSecretRepository) RemoveSecretsExpired() (int64, error) {
	return 0, nil
}

func (repository *MockSecretRepository) HasSecretWithCustomPwd(id string) (bool, error) {
	return false, nil
}
