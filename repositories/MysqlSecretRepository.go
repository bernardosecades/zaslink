package repositories

import (
	"github.com/bernardosecades/sharesecret/models"
)

// https://blog.gojekengineering.com/the-many-flavours-of-dependency-injection-in-go-25aa070d79a0
type MySqlSecretRepository struct {
}

func NewMySqlSecretRepository() SecretRepository {
	return &MySqlSecretRepository{}
}

// https://github.com/irahardianto/service-pattern-go/blob/master/services/PlayerService.go
func (repository *MySqlSecretRepository) GetSecret(id string) (models.Secret, error) {
	return models.Secret{}, nil
}
