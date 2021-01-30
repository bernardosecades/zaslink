package repository

import (
	"github.com/bernardosecades/sharesecret/types"
	"time"
)

type SecretRepository interface {
	GetSecret(id string) (types.Secret, error)
	CreateSecret(content string, customPwd bool, expire time.Time) (types.Secret, error)
	RemoveSecret(id string) error
	HasSecretWithCustomPwd(id string) (bool, error)
}
