package main

import (
	"context"
	"fmt"
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
	a.Telemetry.LogApplicationStart()

	// Wait for the context to be cancelled
	select {
	case <-ctx.Done():
		// Context was cancelled, indicating a shutdown signal
		fmt.Println("Shutting down application...")

		// Log application shutdown
		a.Telemetry.LogApplicationShutdown(ctx)
		defer a.Telemetry.Shutdown(ctx)

		return true
	}
}
