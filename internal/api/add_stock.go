package api

import (
	"context"

	"github.com/davidoram/beaker/internal/db"
	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/internal/utility"
	"github.com/davidoram/beaker/schemas"
	"github.com/nats-io/nats.go/micro"
)

func (app *App) stockAddHandler(ctx context.Context, req micro.Request) {
	rs := newRequestScope(ctx, req, app.nc, app.db)
	defer rs.close(ctx)
	rs.validateJSON(ctx, app.compiler, req.Data(), schemas.StockAddRequestSchema)
	stockReq := decodeRequest[schemas.StockAddRequest](ctx, rs)
	resp := rs.makeStockAddResponse(ctx, rs.addStock(ctx, stockReq))
	rs.commitOrRollback(ctx)
	rs.respondJSON(ctx, req, resp)
}

// addStock adds stock to the inventory.
func (rs *requestScope) addStock(ctx context.Context, req schemas.StockAddRequest) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "add stock")
	defer span.End()

	if rs.hasError() {
		return nil
	}

	params := db.AddInventoryParams{
		ProductSku: string(req.ProductSKU),
		StockLevel: int32(req.Quantity),
	}
	inventory, err := rs.queries.AddInventory(ctx, params)
	if err != nil {
		rs.addCallerError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *requestScope) makeStockAddResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockAddResponse {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "build stock-add response")
	defer span.End()

	resp := schemas.StockAddResponse{}
	if rs.hasError() {
		resp.OK = false
		resp.Error = utility.Ptr(rs.getError().Error())
	} else {
		resp.OK = true
		resp.ProductSKU = utility.Ptr(inventory.ProductSku)
		resp.Quantity = utility.Ptr(int(inventory.StockLevel))
	}
	return &resp
}
