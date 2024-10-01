package middleware

import (
	"net/http"
	"time"

	"github.com/bernardosecades/sharesecret/pkg/api"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewMetricMiddleware creates the middleware that will record all
// HTTP-related metrics.
func NewMetricMiddleware(meter metric.Meter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		durationHistogram, _ := meter.Int64Histogram(
			"http.server.latency",
			metric.WithDescription("latency ms each endpoint"),
			metric.WithUnit("ms"),
		)
		counter, _ := meter.Int64Counter(
			"request_count",
			metric.WithDescription("Incoming request count"),
			metric.WithUnit("request"),
		)

		return &httpMetricMiddleware{
			next:                     next,
			requestDurationHistogram: durationHistogram,
			requestCounter:           counter,
		}
	}
}

// httpMetricMiddleware executes the HTTP endpoint while keeping track
// of how much time it took to execute and add some extra routing information
// to all metrics
type httpMetricMiddleware struct {
	next                     http.Handler
	requestDurationHistogram metric.Int64Histogram
	requestCounter           metric.Int64Counter
}

func (m *httpMetricMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Golang's http implementation doesn't allow us to retrieve the raw response
	// so we have to use a capturer in order to have access to any response
	// data. In this case, we want to capture the response HTTP status code.
	rw := api.NewStatusCodeCapturerWriter(w)

	initialTime := time.Now()
	m.next.ServeHTTP(rw, r)
	duration := time.Since(initialTime)

	route := mux.CurrentRoute(r)
	pathTemplate, _ := route.GetPathTemplate()

	metricAttributes := attribute.NewSet(
		attribute.KeyValue{
			Key:   semconv.HTTPRouteKey,
			Value: attribute.StringValue(pathTemplate),
		},
		attribute.KeyValue{
			Key:   semconv.HTTPRequestMethodKey,
			Value: attribute.StringValue(r.Method),
		},
		attribute.KeyValue{
			Key:   semconv.HTTPResponseStatusCodeKey,
			Value: attribute.IntValue(rw.StatusCode),
		},
	)

	m.requestDurationHistogram.Record(
		r.Context(),
		duration.Milliseconds(),
		metric.WithAttributeSet(metricAttributes),
	)

	m.requestCounter.Add(r.Context(), 1, metric.WithAttributeSet(metricAttributes))

}
