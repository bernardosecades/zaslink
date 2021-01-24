package models

type Secret struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Viewed  bool   `json:"viewed"`
}
