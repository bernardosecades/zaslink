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
		"1328e756e51bced3e067876edc75a65c2aecf06ddf5b7b9e2738d4ab410f5802576488c22ab535",
	}, nil
}
