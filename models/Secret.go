package models

type Secret struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}
// TODO hex.encode, hex.decode