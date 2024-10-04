package middleware

import (
	"net/http"
	"time"

	"github.com/bernardosecades/sharesecret/pkg/api"

	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Golang's http implementation doesn't allow us to retrieve the raw response
			// so we have to use a capturer in order to have access to any response
			// data. In this case, we want to capture the response HTTP status code.
			rw := api.NewStatusCodeCapturerWriter(w)

			start := time.Now().UTC()
			defer func() {
				logger.Info().
					Str("request-id", GetRequestID(r.Context())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("ip", r.RemoteAddr).
					Int("status", rw.StatusCode).
					Int64("latency_ms", time.Since(start).Milliseconds()).
					Str("user-agent", r.UserAgent()).
					Msg("request completed")
			}()
			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}
