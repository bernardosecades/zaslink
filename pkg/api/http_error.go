package api

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Detail     string `json:"detail"`
	StatusCode int    `json:"statusCode"`
}

func NewHTTPError(detail string, statusCode int) *HTTPError {
	return &HTTPError{Detail: detail, StatusCode: statusCode}
}

func (e *HTTPError) Error() string {
	return e.Detail
}

func EncodeHTTPError(err *HTTPError, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(err)
}
