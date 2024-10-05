package events

import (
	"github.com/bernardosecades/sharesecret/internal/entity"
	"github.com/bernardosecades/sharesecret/pkg/events"
)

func NewSecretExpired(secret entity.Secret) events.Event[map[string]string] {
	return events.Event[map[string]string]{
		Name: "secret.expired",
		Data: secret.ToMap(),
	}
}
