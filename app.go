package main

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
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

	for i := 0; i < 100; i++ {

		// Test span to check if telemetry is working
		ctx, span := Tracer.Start(ctx, "application startup", trace.WithSpanKind(trace.SpanKindServer))
		slog.InfoContext(ctx, "Application started")
		span.End()
	}
	println("Application started successfully")

	// Wait for the context to be cancelled
	select {
	case <-ctx.Done():
		println("Shutdown signal received, exiting application...")
		// Context was cancelled, indicating a shutdown signal
		slog.InfoContext(ctx, "Shutting down application...")
		return true
	}
}
