package dto

type ErrorResponse struct {
	StatusCode int    `json:"StatusCode"`
	Err        string `json:"Error"`
}
