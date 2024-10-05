package events

import (
	"github.com/bernardosecades/sharesecret/internal/entity"
	"github.com/bernardosecades/sharesecret/pkg/events"
)

func NewSecretDeleted(secret entity.Secret) events.Event[map[string]string] {
	return events.Event[map[string]string]{
		Name: "secret.deleted",
		Data: secret.ToMap(),
	}
}
