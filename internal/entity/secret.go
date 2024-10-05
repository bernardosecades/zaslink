package entity

import (
	"strconv"
	"time"
)

type Secret struct {
	ID        string    `bson:"_id"`
	PrivateID string    `bson:"privateId"`
	Content   []byte    `bson:"content"`
	CustomPwd bool      `bson:"customPwd"`
	Viewed    bool      `bson:"viewed"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updateAt"`
	ExpiredAt time.Time `bson:"expiredAt"`
}

func (s *Secret) ToMap() map[string]string {
	return map[string]string{
		"id":        s.ID,
		"privateId": s.PrivateID,
		"createdAt": s.CreatedAt.String(),
		"updatedAt": s.UpdatedAt.String(),
		"expiredAt": s.ExpiredAt.String(),
		"customPwd": strconv.FormatBool(s.CustomPwd),
	}
}
