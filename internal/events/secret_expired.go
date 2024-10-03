package events

import (
	"strconv"
	"time"

	"github.com/bernardosecades/sharesecret/pkg/events"

	"github.com/bernardosecades/sharesecret/internal/entity"
)

func NewSecretExpired(secret entity.Secret) events.Event[map[string]string] {
	return events.Event[map[string]string]{
		Name: "secret.expired",
		Data: map[string]string{
			"id":        secret.ID,
			"createdAt": secret.CreatedAt.Format(time.RFC3339),
			"expiredAt": secret.ExpiredAt.Format(time.RFC3339),
			"customPwd": strconv.FormatBool(secret.CustomPwd),
		},
	}
}
