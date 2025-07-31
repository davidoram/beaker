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
		Name:        "InventoryService",
		Version:     "0.1.0",
		Description: "Manage product inventory",
	}
	svc, err := micro.AddService(app.nc, config)
	if err != nil {
		return err
	}
	// add a group to aggregate endpoints under common prefix
	inventory := svc.AddGroup("inventory")
	err = inventory.AddEndpoint("receive", micro.HandlerFunc(app.inventoryReceiveHandler))
	if err != nil {
		return err
	}
	err = inventory.AddEndpoint("drawdown", micro.HandlerFunc(app.inventoryDrawdownHandler))
	if err != nil {
		return err
	}
	err = inventory.AddEndpoint("show", micro.HandlerFunc(app.inventoryShowHandler))
	if err != nil {
		return err
	}
	app.svc = svc
	return nil
}

func (app *App) inventoryReceiveHandler(req micro.Request) {
	req.Respond([]byte("TODO"))
}

func (app *App) inventoryDrawdownHandler(req micro.Request) {
	req.Respond([]byte("TODO"))
}

func (app *App) inventoryShowHandler(req micro.Request) {
	req.Respond([]byte("TODO"))
}
