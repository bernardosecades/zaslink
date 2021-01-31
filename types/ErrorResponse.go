package types

type ErrorResponse struct {
	StatusCode int    `json:"Status"`
	Err        string `json:"Error"`
}
