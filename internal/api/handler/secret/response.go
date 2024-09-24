package secret

import (
	"time"
)

type CreateSecretResponse struct {
	ID        string    `json:"id"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type RetrieveSecretResponse struct {
	Content string `json:"content"`
}
