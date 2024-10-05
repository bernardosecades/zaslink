package events

import (
	"github.com/bernardosecades/sharesecret/internal/entity"
	"github.com/bernardosecades/sharesecret/pkg/events"
)

func NewSecretCreated(secret entity.Secret) events.Event[map[string]string] {
	return events.Event[map[string]string]{
		Name: "secret.created",
		Data: secret.ToMap(),
	}
}
