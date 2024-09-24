package entity

import "time"

type Secret struct {
	ID        string    `bson:"_id"`
	Content   string    `bson:"content"`
	CustomPwd bool      `bson:"customPwd"`
	Viewed    bool      `bson:"viewed"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiredAt time.Time `bson:"expiredAt"`
}
