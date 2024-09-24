package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	requestIDKey    ctxKey = ctxKey(0)
	requestIDHeader        = "X-Request-Id"
)

func GetRequestID(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		w.Header().Set(requestIDHeader, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
