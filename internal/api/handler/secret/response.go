package secret

import (
	"encoding/json"
	"fmt"
	"time"
)

type CreateSecretResponse struct {
	ID        string    `json:"id"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type RetrieveSecretResponse struct {
	Content string `json:"content"`
}

type HTTPErrorResponse struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}

func (e *HTTPErrorResponse) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

// ResponseHeaders returns http status code and headers.
func (e *HTTPErrorResponse) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

// ResponseBody returns JSON response body.
func (e *HTTPErrorResponse) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("error while parsing response body: %v", err)
	}
	return body, nil
}

func NewHTTPError(err error, status int, detail string) error {
	return &HTTPErrorResponse{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}
