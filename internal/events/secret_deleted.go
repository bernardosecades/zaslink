package events

import (
	"github.com/bernardosecades/zaslink/internal/entity"
	"github.com/bernardosecades/zaslink/pkg/events"
)

func NewSecretDeleted(secret entity.Secret) events.Event[map[string]string] {
	return events.Event[map[string]string]{
		Name: "secret.deleted",
		Data: secret.ToMap(),
	}
}
