package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

// TelemetryProvider manages the OpenTelemetry tracing provider
type TelemetryProvider struct {
	tracerProvider *trace.TracerProvider
	serviceName    string
}

// NewTelemetryProvider creates and configures an OpenTelemetry TracerProvider
func NewTelemetryProvider(serviceName string) (*TelemetryProvider, error) {
	// Create a resource describing this application
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// Standard service resource attributes
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("0.1.0"),
			attribute.String("environment", getEnv("OTEL_ENVIRONMENT", "development")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create the trace exporter
	exporter, err := createTraceExporter(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create a tracer provider with the exporter
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	return &TelemetryProvider{
		tracerProvider: tp,
		serviceName:    serviceName,
	}, nil
}

// createTraceExporter creates an appropriate exporter based on environment configuration
func createTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	// Default to OTLP exporter connecting to the local collector
	endpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")

	// Check if we should use OTLP or stdout exporter
	exporterType := getEnv("OTEL_EXPORTER_TYPE", "otlp")

	if exporterType == "stdout" {
		// Use stdout exporter for local development or debugging
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	}

	// Use OTLP gRPC exporter for production use
	return otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(), // For dev environment; configure TLS for production
	)
}

// Shutdown cleanly shuts down the telemetry provider
func (tp *TelemetryProvider) Shutdown(ctx context.Context) error {
	// Create a timeout for shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Shutdown the tracer provider
	if err := tp.tracerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}
	return nil
}

// LogApplicationStart logs a standard application startup event
func (tp *TelemetryProvider) LogApplicationStart() {
	tracer := otel.Tracer(tp.serviceName)

	// Create a new context to start the span
	ctx := context.Background()

	// Create a span for application startup
	_, span := tracer.Start(ctx, "application_start")
	span.SetAttributes(
		attribute.String("event", "application_start"),
		attribute.String("service.name", tp.serviceName),
		attribute.String("timestamp", time.Now().Format(time.RFC3339)),
	)
	span.End()
}

// LogApplicationShutdown logs a standard application shutdown event
func (tp *TelemetryProvider) LogApplicationShutdown(ctx context.Context) {
	tracer := otel.Tracer(tp.serviceName)

	// Create a span for application shutdown
	_, span := tracer.Start(ctx, "application_shutdown")
	span.SetAttributes(
		attribute.String("event", "application_shutdown"),
		attribute.String("service.name", tp.serviceName),
		attribute.String("timestamp", time.Now().Format(time.RFC3339)),
	)
	span.End()
}

// Helper function to get environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
