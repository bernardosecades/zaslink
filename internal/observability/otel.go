package observability

import (
	"time"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewTraceExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

// NewMetricExporter output in stdout only keep to debug
func NewMetricExporter() (metric.Exporter, error) {
	return stdoutmetric.New()
}

func NewTraceProvider(traceExporter trace.SpanExporter) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)),
	)
	return traceProvider
}

// TODO merge with main.go provider
func NewMeterProvider(meterExporter metric.Exporter) *metric.MeterProvider {
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(meterExporter,
			metric.WithInterval(10*time.Second))), // TODO DEFAULT ONE: 1 MINUTE
	)
	return meterProvider
}

// NewPropagator For trace information to be sent to remote processes, we need to propagate the context
func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)
}
