package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("instrumentation/beaker/main")
}

// GetTracer returns the global tracer instance
func GetTracer() trace.Tracer {
	return tracer
}
