package repositories

import (
	"github.com/bernardosecades/sharesecret/models"
	"time"
)

type SecretRepository interface {
	GetSecret(id string) (models.Secret, error)
	CreateSecret(content string, customPwd bool, expire time.Time) (models.Secret, error)
	RemoveSecret(id string) error
	HasSecretWithCustomPwd(id string) (bool, error)
}
