package repositories

import (
	"github.com/bernardosecades/sharesecret/models"
)

// https://blog.gojekengineering.com/the-many-flavours-of-dependency-injection-in-go-25aa070d79a0
type InMemorySecret struct {
}

func NewInMemorySecretRepository() SecretRepository {
	return &InMemorySecret{}
}

func (r *InMemorySecret) GetSecret(id string) (models.Secret, error) {
	return models.Secret{
		"727d7040-aac7-4dc3-ab44-938bfba92ebd",
		"cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922",
	}, nil
}
