package api

import (
	"encoding/json"
	"net/http"
)

func EncodeHTTPResponse(v interface{}, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
