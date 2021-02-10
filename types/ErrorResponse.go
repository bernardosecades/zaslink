package types

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Err        string `json:"error"`
}
