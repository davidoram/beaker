package main

import (
	"context"
	"log/slog"
)

type App struct {
	Options   Options
	Telemetry *TelemetryProvider
}

func NewApp(opts Options, telemetry *TelemetryProvider) *App {

	return &App{
		Options:   opts,
		Telemetry: telemetry,
	}
}

func (a *App) Start(ctx context.Context) bool {

	_, span := a.Telemetry.tracerProvider.Tracer("").Start(ctx, "application_start")
	span.End()

	slog.InfoContext(ctx, "Application started")
	// Wait for the context to be cancelled
	select {
	case <-ctx.Done():
		// Context was cancelled, indicating a shutdown signal
		slog.InfoContext(ctx, "Shutting down application...")
		return true
	}
}
