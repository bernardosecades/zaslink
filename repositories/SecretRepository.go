package repositories

import (
	"github.com/bernardosecades/sharesecret/models"
)

type SecretRepository interface {
	GetSecret(id string) (models.Secret, error)
}
