package events

import (
	"github.com/bernardosecades/zaslink/internal/entity"
	"github.com/bernardosecades/zaslink/pkg/events"
)

func NewSecretViewed(secret entity.Secret) events.Event[map[string]string] {
	return events.Event[map[string]string]{
		Name: "secret.viewed",
		Data: secret.ToMap(),
	}
}
