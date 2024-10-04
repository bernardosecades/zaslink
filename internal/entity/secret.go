package entity

import "time"

type Secret struct {
	ID        string    `bson:"_id"`
	Content   []byte    `bson:"content"`
	CustomPwd bool      `bson:"customPwd"`
	Viewed    bool      `bson:"viewed"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updateAt"`
	ExpiredAt time.Time `bson:"expiredAt"`
}
