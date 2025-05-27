package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

// TelemetryProvider manages the OpenTelemetry providers
type TelemetryProvider struct {
	tracerProvider *sdktrace.TracerProvider
	serviceName    string
	meter          metric.Meter
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
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Configure slog to send logs to OpenTelemetry
	handler := otelslog.NewHandler(serviceName, otelslog.WithSource(true))
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Get a meter for metrics
	meter := otel.Meter(serviceName)

	return &TelemetryProvider{
		tracerProvider: tp,
		serviceName:    serviceName,
		meter:          meter,
	}, nil
}

// createTraceExporter creates an appropriate exporter based on environment configuration
func createTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
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

// // LogApplicationStart logs a standard application startup event
// func (tp *TelemetryProvider) LogApplicationStart(ctx context.Context) {
// 	// Log the application start event
// 	slog.InfoContext(ctx, "Application started",
// 		slog.String("event", "application_start"),
// 		slog.String("service.name", tp.serviceName),
// 		slog.Time("timestamp", time.Now()),
// 	)

// 	// Create a span for application startup
// 	tracer := otel.Tracer(tp.serviceName)
// 	_, span := tracer.Start(ctx, "application_start")
// 	span.SetAttributes(
// 		attribute.String("event", "application_start"),
// 		attribute.String("service.name", tp.serviceName),
// 	)
// 	span.End()

// 	// Record a metric for application start
// 	counter, _ := tp.meter.Int64Counter("application.starts")
// 	counter.Add(ctx, 1, attribute.String("service.name", tp.serviceName))
// }

// // LogApplicationShutdown logs a standard application shutdown event
// func (tp *TelemetryProvider) LogApplicationShutdown(ctx context.Context) {
// 	// Log the application shutdown event
// 	slog.InfoContext(ctx, "Application shutting down",
// 		slog.String("event", "application_shutdown"),
// 		slog.String("service.name", tp.serviceName),
// 		slog.Time("timestamp", time.Now()),
// 	)

// 	// Create a span for application shutdown
// 	tracer := otel.Tracer(tp.serviceName)
// 	_, span := tracer.Start(ctx, "application_shutdown")
// 	span.SetAttributes(
// 		attribute.String("event", "application_shutdown"),
// 		attribute.String("service.name", tp.serviceName),
// 	)
// 	span.End()

// 	// Record a metric for application shutdown
// 	counter, _ := tp.meter.Int64Counter("application.shutdowns")
// 	counter.Add(ctx, 1, attribute.String("service.name", tp.serviceName))
// }

// Helper function to get environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
