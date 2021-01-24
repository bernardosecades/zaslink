package models

type Secret struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Visited bool   `json:"visited"`
}
