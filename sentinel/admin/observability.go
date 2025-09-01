package admin

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// ObservabilityManager handles metrics, traces, and logging
type ObservabilityManager struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	tracer         trace.Tracer
	meter          metric.Meter

	// Metrics
	requestCounter       metric.Int64Counter
	errorCounter         metric.Int64Counter
	latencyHistogram     metric.Float64Histogram
	violationCounter     metric.Int64Counter
	cryptoOperationTimer metric.Float64Histogram
	detectionTimer       metric.Float64Histogram
	redactionTimer       metric.Float64Histogram
}

// NewObservabilityManager creates a new observability manager
func NewObservabilityManager() (*ObservabilityManager, error) {
	// Create stdout exporters for metrics and traces
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)))
	otel.SetMeterProvider(meterProvider)

	// Create tracer provider
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExporter))
	otel.SetTracerProvider(tracerProvider)

	// Get tracer and meter
	tracer := tracerProvider.Tracer("sentinel-tracer")
	meter := meterProvider.Meter("sentinel-meter")

	// Create metrics
	requestCounter, err := meter.Int64Counter("sentinel.requests.total", metric.WithDescription("Total number of requests"))
	if err != nil {
		return nil, fmt.Errorf("failed to create request counter: %w", err)
	}

	errorCounter, err := meter.Int64Counter("sentinel.errors.total", metric.WithDescription("Total number of errors"))
	if err != nil {
		return nil, fmt.Errorf("failed to create error counter: %w", err)
	}

	latencyHistogram, err := meter.Float64Histogram("sentinel.request.latency", metric.WithDescription("Request latency in seconds"), metric.WithUnit("s"))
	if err != nil {
		return nil, fmt.Errorf("failed to create latency histogram: %w", err)
	}

	violationCounter, err := meter.Int64Counter("sentinel.violations.total", metric.WithDescription("Total number of security violations detected"))
	if err != nil {
		return nil, fmt.Errorf("failed to create violation counter: %w", err)
	}

	cryptoOperationTimer, err := meter.Float64Histogram("sentinel.crypto.operation.duration", metric.WithDescription("Duration of crypto operations"), metric.WithUnit("s"))
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto operation timer: %w", err)
	}

	detectionTimer, err := meter.Float64Histogram("sentinel.detection.duration", metric.WithDescription("Duration of detection operations"), metric.WithUnit("s"))
	if err != nil {
		return nil, fmt.Errorf("failed to create detection timer: %w", err)
	}

	redactionTimer, err := meter.Float64Histogram("sentinel.redaction.duration", metric.WithDescription("Duration of redaction operations"), metric.WithUnit("s"))
	if err != nil {
		return nil, fmt.Errorf("failed to create redaction timer: %w", err)
	}

	return &ObservabilityManager{
		tracerProvider:       tracerProvider,
		meterProvider:        meterProvider,
		tracer:               tracer,
		meter:                meter,
		requestCounter:       requestCounter,
		errorCounter:         errorCounter,
		latencyHistogram:     latencyHistogram,
		violationCounter:     violationCounter,
		cryptoOperationTimer: cryptoOperationTimer,
		detectionTimer:       detectionTimer,
		redactionTimer:       redactionTimer,
	}, nil
}

// StartTrace starts a new trace span
func (o *ObservabilityManager) StartTrace(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return o.tracer.Start(ctx, name, opts...)
}

// RecordMetric records a metric with attributes
func (o *ObservabilityManager) RecordMetric(ctx context.Context, name string, value float64, attrs ...attribute.KeyValue) {
	// This is a simplified implementation - in practice, you'd use specific metric instruments
	switch name {
	case "request.count":
		o.requestCounter.Add(ctx, int64(value), metric.WithAttributes(attrs...))
	case "error.count":
		o.errorCounter.Add(ctx, int64(value), metric.WithAttributes(attrs...))
	case "latency":
		o.latencyHistogram.Record(ctx, value, metric.WithAttributes(attrs...))
	case "violation.count":
		o.violationCounter.Add(ctx, int64(value), metric.WithAttributes(attrs...))
	case "crypto.duration":
		o.cryptoOperationTimer.Record(ctx, value, metric.WithAttributes(attrs...))
	case "detection.duration":
		o.detectionTimer.Record(ctx, value, metric.WithAttributes(attrs...))
	case "redaction.duration":
		o.redactionTimer.Record(ctx, value, metric.WithAttributes(attrs...))
	}
}

// LogEvent logs an event with structured data
func (o *ObservabilityManager) LogEvent(level string, message string, fields map[string]interface{}) {
	// In a real implementation, you'd use a structured logging library
	log.Printf("[%s] %s", level, message)
	for k, v := range fields {
		log.Printf("  %s: %v", k, v)
	}
}

// RecordCryptoOperation records timing for crypto operations
func (o *ObservabilityManager) RecordCryptoOperation(ctx context.Context, operation string, duration time.Duration, attrs ...attribute.KeyValue) {
	additionalAttrs := []attribute.KeyValue{
		attribute.String("operation", operation),
	}
	additionalAttrs = append(additionalAttrs, attrs...)
	o.cryptoOperationTimer.Record(ctx, duration.Seconds(), metric.WithAttributes(additionalAttrs...))
}

// RecordDetectionOperation records timing for detection operations
func (o *ObservabilityManager) RecordDetectionOperation(ctx context.Context, detector string, duration time.Duration, attrs ...attribute.KeyValue) {
	additionalAttrs := []attribute.KeyValue{
		attribute.String("detector", detector),
	}
	additionalAttrs = append(additionalAttrs, attrs...)
	o.detectionTimer.Record(ctx, duration.Seconds(), metric.WithAttributes(additionalAttrs...))
}

// RecordRedactionOperation records timing for redaction operations
func (o *ObservabilityManager) RecordRedactionOperation(ctx context.Context, action string, duration time.Duration, attrs ...attribute.KeyValue) {
	additionalAttrs := []attribute.KeyValue{
		attribute.String("action", action),
	}
	additionalAttrs = append(additionalAttrs, attrs...)
	o.redactionTimer.Record(ctx, duration.Seconds(), metric.WithAttributes(additionalAttrs...))
}

// RecordViolation records a security violation
func (o *ObservabilityManager) RecordViolation(ctx context.Context, violationType string, attrs ...attribute.KeyValue) {
	additionalAttrs := []attribute.KeyValue{
		attribute.String("violation_type", violationType),
	}
	additionalAttrs = append(additionalAttrs, attrs...)
	o.violationCounter.Add(ctx, 1, metric.WithAttributes(additionalAttrs...))
}

// Shutdown cleans up resources
func (o *ObservabilityManager) Shutdown(ctx context.Context) error {
	if err := o.tracerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}
	if err := o.meterProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown meter provider: %w", err)
	}
	return nil
}
