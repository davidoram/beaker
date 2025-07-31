package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type App struct {
	nc  *nats.Conn
	db  *pgxpool.Pool
	svc micro.Service
}

func StartNewApp(nc *nats.Conn, db *pgxpool.Pool) (*App, error) {
	app := &App{
		nc: nc,
		db: db,
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
	err = stock.AddEndpoint("add", micro.HandlerFunc(app.stockAddHandler))
	if err != nil {
		return err
	}
	err = stock.AddEndpoint("remove", micro.HandlerFunc(app.stockRemoveHandler))
	if err != nil {
		return err
	}
	err = stock.AddEndpoint("get", micro.HandlerFunc(app.stockGetHandler))
	if err != nil {
		return err
	}
	app.svc = svc
	return nil
}

func (app *App) stockAddHandler(req micro.Request) {
	req.Respond([]byte("TODO"))
}

func (app *App) stockRemoveHandler(req micro.Request) {
	req.Respond([]byte("TODO"))
}

func (app *App) stockGetHandler(req micro.Request) {
	req.Respond([]byte("TODO"))
}
