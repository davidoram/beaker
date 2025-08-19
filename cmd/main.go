package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davidoram/beaker/internal/api"
	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/internal/utility"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"golang.org/x/exp/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setTimeZoneToUTCOrExit()
	opts := getOptionsOrExit()
	closeOtel := setupTelemetryOrExit(ctx, opts)
	defer closeOtel()
	pool := setupPostgresPoolOrExit(ctx, opts.PostgresURL)
	defer pool.Close()
	nc := connectToNATSOrExit(ctx, opts.NatsURL, opts.CredentialsFile)
	defer nc.Close()
	compiler := makeJSONSchemaCompilerOrExit(ctx, opts.SchemaDir)
	setupSignalHandler(ctx, cancel)
	app := startAppOrExit(nc, pool, compiler)
	defer app.Stop()
	slog.InfoContext(ctx, "beaker is running")

	// Wait for the context to be cancelled
	<-ctx.Done()
	slog.InfoContext(ctx, "shutting down application")
}

func getOptionsOrExit() Options {
	// Parse command line arguments to get options
	opts, err := GetOptions()
	if err != nil {
		// If options parsing fails, print the error and exit with an error code
		os.Stderr.WriteString("Error parsing options: " + err.Error() + "\n")
		os.Exit(1)
	}
	return opts
}

func setTimeZoneToUTCOrExit() {
	// Set the timezone to UTC
	if err := os.Setenv("TZ", "UTC"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set timezone to UTC: %v\n", err)
		os.Exit(1)
	}
}

func setupTelemetryOrExit(ctx context.Context, opts Options) func() {
	// Initialize telemetry
	shutdown, err := telemetry.NewTelemetryProvider(ctx, "beaker")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize telemetry: %v\n", err)
		os.Exit(1)
	}
	return func() { shutdown(ctx) }
}

func setupPostgresPoolOrExit(ctx context.Context, postgresUrl string) *pgxpool.Pool {
	// Optional: configure pool settings
	config, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to parse db config", err)
		os.Exit(1)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute

	// Set up the tracer for OpenTelemetry
	config.ConnConfig.Tracer = otelpgx.NewTracer()

	// Connect and create the pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to create pool", err)
		os.Exit(1)
	}

	if err := otelpgx.RecordStats(pool); err != nil {
		slog.ErrorContext(ctx, "unable to record database stats", err)
		os.Exit(1)
	}

	// Ping the database to ensure it's reachable
	if err := pool.Ping(ctx); err != nil {
		slog.ErrorContext(ctx, "Unable to ping database", err)
		os.Exit(1)
	}
	slog.InfoContext(ctx, "postgres connection pool created successfully")
	return pool
}

func connectToNATSOrExit(ctx context.Context, natsURL, credentialsFile string) *nats.Conn {
	// Connect to NATS server
	nc, err := nats.Connect(
		natsURL,
		nats.Name("beaker"),
		nats.UserCredentials(credentialsFile),
		nats.MaxReconnects(-1),            // Infinite reconnect attempts
		nats.ReconnectWait(2*time.Second), // Wait 2 seconds between reconnect attempts
	)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to connect to NATS server", err)
		os.Exit(1)
	}
	return nc
}

func makeJSONSchemaCompilerOrExit(ctx context.Context, schemaDir string) *jsonschema.Compiler {
	loader, err := utility.NewLoader(map[string]string{
		"http://github.com/davidoram/beaker/schemas/": schemaDir,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Unable to create JSON schema loader", err)
		os.Exit(1)
	}
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(loader)
	compiler.AssertContent()
	compiler.AssertFormat()
	compiler.DefaultDraft(jsonschema.Draft2020)
	return compiler
}

func setupSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh // Receive signal but don't need to use it
		slog.InfoContext(ctx, "Received termination signal, shutting down...")
		cancel() // Cancel the context to signal graceful shutdown
	}()
}

func startAppOrExit(nc *nats.Conn, pool *pgxpool.Pool, compiler *jsonschema.Compiler) *api.App {
	// Create a new application instance
	app, err := api.StartNewApp(nc, pool, compiler)
	if err != nil {
		slog.Error("Failed to create application instance", err)
		os.Exit(1)
	}
	return app
}
