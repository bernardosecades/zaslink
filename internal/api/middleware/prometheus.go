package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type ResponseWriter struct {
	http.ResponseWriter
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

		next.ServeHTTP(w, r)

		TotalRequests.WithLabelValues(path).Inc()
	})
}
