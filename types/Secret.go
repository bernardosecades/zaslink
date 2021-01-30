package types

import "time"

type Secret struct {
	Id        string `json:"id"`
	Content   string `json:"content"`
	CustomPwd bool   `json:"customPwd"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}
