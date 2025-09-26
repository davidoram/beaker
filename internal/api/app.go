package api

import (
	"context"
	"log/slog"

	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// App represents the application context
type App struct {
	nc       *nats.Conn
	db       *pgxpool.Pool
	svc      micro.Service
	compiler *jsonschema.Compiler
}

func StartNewApp(nc *nats.Conn, db *pgxpool.Pool, compiler *jsonschema.Compiler) (*App, error) {

	app := &App{
		nc:       nc,
		db:       db,
		compiler: compiler,
	}
	if err := app.makeService(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Stop() error {
	if err := app.svc.Stop(); err != nil {
		return err
	}
	return nil
}

func (app *App) makeService() error {
	config := micro.Config{
		Name:        "StockService",
		Version:     "0.1.0",
		Description: "Manage product stock",
		ErrorHandler: func(svc micro.Service, err *micro.NATSError) {
			slog.Error("Service error occurred", "error", err)
		},
	}
	svc, err := micro.AddService(app.nc, config)
	if err != nil {
		return err
	}
	// add a group to aggregate endpoints under common prefix
	stock := svc.AddGroup("stock")
	err = stock.AddEndpoint("add", micro.HandlerFunc(traceHandler(app.stockAddHandler)))
	if err != nil {
		return err
	}
	err = stock.AddEndpoint("remove", micro.HandlerFunc(traceHandler(app.stockRemoveHandler)))
	if err != nil {
		return err
	}
	err = stock.AddEndpoint("get", micro.HandlerFunc(traceHandler(app.stockGetHandler)))
	if err != nil {
		return err
	}
	app.svc = svc
	return nil
}

func traceHandler(handler func(ctx context.Context, req micro.Request)) micro.HandlerFunc {
	return func(req micro.Request) {
		ctx := context.Background()
		// Start a new otel trace span
		tracer := telemetry.GetTracer()
		ctx, span := tracer.Start(ctx, req.Subject())
		defer span.End()
		slog.InfoContext(ctx, "API Request "+req.Subject())
		handler(ctx, req)
	}
}
