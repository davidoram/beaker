package api

import (
	"context"
	"log/slog"

	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

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

func (app *App) stockAddHandler(ctx context.Context, req micro.Request) {
	rs := NewRequestScope(ctx, req, app.db)
	defer rs.Close(ctx)
	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockAddRequestSchema)
	stockReq := DecodeRequest[schemas.StockAddRequest](ctx, rs)
	resp := rs.MakeStockAddResponse(ctx, rs.AddStock(ctx, stockReq))
	rs.CommitOrRollback(ctx)
	rs.RespondJSON(ctx, req, resp)
}

func (app *App) stockRemoveHandler(ctx context.Context, req micro.Request) {
	rs := NewRequestScope(ctx, req, app.db)
	defer rs.Close(ctx)
	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockRemoveRequestSchema)
	stockReq := DecodeRequest[schemas.StockRemoveRequest](ctx, rs)
	resp := rs.MakeStockRemoveResponse(ctx, rs.RemoveStock(ctx, stockReq))
	rs.CommitOrRollback(ctx)
	rs.RespondJSON(ctx, req, resp)
}

func (app *App) stockGetHandler(ctx context.Context, req micro.Request) {
	rs := NewRequestScope(ctx, req, app.db)
	defer rs.Close(ctx)
	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockGetRequestSchema)
	stockReq := DecodeRequest[schemas.StockGetRequest](ctx, rs)
	resp := rs.MakeStockGetResponse(ctx, rs.GetStock(ctx, stockReq))
	rs.CommitOrRollback(ctx)
	rs.RespondJSON(ctx, req, resp)
}
