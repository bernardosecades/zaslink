package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var TotalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total_bernie",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

func Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		TotalRequests.WithLabelValues(path).Inc()
	})
}
