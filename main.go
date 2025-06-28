package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"
)

func main() {
	// Parse command line arguments to get options
	opts, err := GetOptions()
	if err != nil {
		// If options parsing fails, print the error and exit with an error code
		os.Stderr.WriteString("Error parsing options: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Create a context for the application and
	// make it cancelled if a termination signal is received
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is cancelled when main exits

	// Initialize telemetry
	shutdown, err := NewTelemetryProvider(ctx, "beaker")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize telemetry: %v\n", err)
		os.Exit(1)
	}
	defer shutdown(ctx) // Ensure telemetry is shut down when main exits

	// Log application start
	slog.InfoContext(ctx, "Telemetry initialized")
	// telemetry.LogApplicationStart(ctx)

	// Initialize the application
	app := NewApp(opts)

	// Add a signal handler to gracefully handle termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh // Receive signal but don't need to use it
		slog.InfoContext(ctx, "Received termination signal, shutting down...")
		cancel() // Cancel the context to signal graceful shutdown
	}()

	// Start the application
	if app.Start(ctx) {
		os.Exit(0)
	}
	os.Exit(1)
}
