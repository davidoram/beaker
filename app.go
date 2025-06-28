package main

import (
	"context"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	name = "beaker"
)

var (
	tracer    = otel.Tracer(name)
	meter     = otel.Meter(name)
	logger    = otelslog.NewLogger(name)
	callCount metric.Int64Counter
)

type App struct {
	Options Options
}

func NewApp(opts Options) *App {

	return &App{
		Options: opts,
	}
}

func (a *App) Start(ctx context.Context) bool {

	// Test span to check if telemetry is working
	ctx, span := tracer.Start(ctx, "application startup", trace.WithSpanKind(trace.SpanKindServer))
	logger.InfoContext(ctx, "Application started")
	span.End()

	// Wait for the context to be cancelled
	select {
	case <-ctx.Done():
		logger.InfoContext(ctx, "Shutting down application...")
		return true
	}
}
