package main

import (
	"context"

	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.opentelemetry.io/otel"
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
	err = stock.AddEndpoint("add", micro.HandlerFunc(tracer(app.stockAddHandler)))
	if err != nil {
		return err
	}
	err = stock.AddEndpoint("remove", micro.HandlerFunc(tracer(app.stockRemoveHandler)))
	if err != nil {
		return err
	}
	err = stock.AddEndpoint("get", micro.HandlerFunc(tracer(app.stockGetHandler)))
	if err != nil {
		return err
	}
	app.svc = svc
	return nil
}

func tracer(handler func(ctx context.Context, req micro.Request)) micro.HandlerFunc {
	return func(req micro.Request) {
		ctx := context.Background()
		// Start a new otel trace span
		ctx, span := otel.Tracer("").Start(ctx, req.Subject())
		defer span.End()
		handler(ctx, req)
	}
}

func (app *App) stockAddHandler(ctx context.Context, req micro.Request) {
	rs, err := NewRequestScope(ctx, req, app.db)
	if err != nil {
		req.RespondJSON(schemas.BuildErrorResponse(&schemas.StockAddResponse{}, err))
		return
	}
	defer rs.Close(ctx)

	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockAddRequestSchema)
	stockReq := DecodeRequest[schemas.StockAddRequest](ctx, rs)
	resp := rs.BuildStockAddResponse(rs.AddStock(ctx, stockReq))
	req.RespondJSON(resp)
}

func (app *App) stockRemoveHandler(ctx context.Context, req micro.Request) {
	req.Respond([]byte("TODO"))
}

func (app *App) stockGetHandler(ctx context.Context, req micro.Request) {
	req.Respond([]byte("TODO"))
}

func validateAgainstSchema(ctx context.Context, rs *requestScope, schema string) {
	// TODO: Implement schema validation logic
}
