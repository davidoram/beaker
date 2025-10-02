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
	rs := NewRequestScope(ctx, req, app.nc, app.db)
	defer rs.Close(ctx)
	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockAddRequestSchema)
	stockReq := DecodeRequest[schemas.StockAddRequest](ctx, rs)
	resp := rs.MakeStockAddResponse(ctx, rs.AddStock(ctx, stockReq))
	rs.CommitOrRollback(ctx)
	rs.RespondJSON(ctx, req, resp)
}

// AddStock adds stock to the inventory.
func (rs *requestScope) AddStock(ctx context.Context, req schemas.StockAddRequest) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "add stock")
	defer span.End()

	if rs.HasError() {
		return nil
	}

	params := db.AddInventoryParams{
		ProductSku: string(req.ProductSKU),
		StockLevel: int32(req.Quantity),
	}
	inventory, err := rs.queries.AddInventory(ctx, params)
	if err != nil {
		rs.AddCallerError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *requestScope) MakeStockAddResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockAddResponse {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "build stock-add response")
	defer span.End()

	resp := schemas.StockAddResponse{}
	if rs.HasError() {
		resp.OK = false
		resp.Error = utility.Ptr(rs.GetError().Error())
	} else {
		resp.OK = true
		resp.ProductSKU = utility.Ptr(inventory.ProductSku)
		resp.Quantity = utility.Ptr(int(inventory.StockLevel))
	}
	return &resp
}
