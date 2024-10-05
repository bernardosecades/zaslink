package secret

type CreateSecretRequest struct {
	Content    string `json:"content"`
	Pwd        string `json:"pwd"`
	Expiration string `json:"expiration"`
}
