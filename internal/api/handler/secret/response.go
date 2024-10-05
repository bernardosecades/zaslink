package secret

import (
	"time"
)

type CreateSecretResponse struct {
	ID        string    `json:"id"`
	PrivateID string    `json:"privateId"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type RetrieveSecretResponse struct {
	Content string `json:"content"`
}
