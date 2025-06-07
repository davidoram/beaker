package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	TelemetryNameSpace = "github.com/davidoram/beaker"
)

var (
	Tracer trace.Tracer
	Meter  *sdkmetric.MeterProvider
	Logger *log.LoggerProvider
)

// NewTelemetryProvider creates and configures an OpenTelemetry TracerProvider
func NewTelemetryProvider(ctx context.Context, serviceName string) (func(context.Context), error) {

	// Build a shutdown function
	var shutdownFuncs []func(context.Context) error
	shutdown := func(ctx context.Context) {
		for _, fn := range shutdownFuncs {
			_ = fn(ctx)
		}
		shutdownFuncs = nil
	}

	// Create a resource describing this application
	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			// Standard service resource attributes
			semconv.ServiceName(getEnv("OTEL_SERVICE_NAME", "beaker")),
			semconv.ServiceVersion("0.1.0"),
			attribute.String("environment", getEnv("OTEL_ENVIRONMENT", "development")),
			semconv.ServiceInstanceID(uuid.NewString()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Set the global propagator to propagate W3C trace context and baggage
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Create a tracer provider that exports traces via GRPC
	// and uses the resource we created
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)
	Tracer = traceProvider.Tracer(serviceName)

	// Set the global meter provider
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	Meter = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithProducer(runtime.NewProducer()),
			),
		),
		sdkmetric.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, Meter.Shutdown)
	otel.SetMeterProvider(Meter)

	// Create a logger provider that uses OpenTelemetry
	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint("localhost:4317"),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	Logger = log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)

	// Create new logger and set it as the default logger
	slog.SetDefault(otelslog.NewLogger(serviceName, otelslog.WithLoggerProvider(Logger)))
	global.SetLoggerProvider(Logger)

	slog.Info("OpenTelemetry setup complete")

	shutdownFuncs = append(shutdownFuncs, Logger.Shutdown)

	err = runtime.Start(runtime.WithMeterProvider(Meter))
	if err != nil {
		return shutdown, err
	}
	return shutdown, nil
}

// Helper function to get environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
