package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				// TODO if response is error print error. Use HttpError to print Unwrap error
				logger.Info().
					Str("request-id", GetRequestID(r.Context())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("ip", r.RemoteAddr).
					Int("status", 200). // TODO make wrapper to get status because by default you can not get it
					Int64("latency", time.Since(start).Milliseconds()).
					Str("user-agent", r.UserAgent()).
					Msg("request completed")
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
