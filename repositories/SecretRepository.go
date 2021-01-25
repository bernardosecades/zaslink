package repositories

import (
	"github.com/bernardosecades/sharesecret/models"
)

type SecretRepository interface {
	GetSecret(id string) (models.Secret, error)
	CreateSecret(content string, customPwd bool) (models.Secret, error)
	UpdateToViewed(id string) error
	HasSecretWithCustomPwd(id string) (bool, error)
}
