package api

import (
	"context"

	"github.com/davidoram/beaker/internal/db"
	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/internal/utility"
	"github.com/davidoram/beaker/schemas"
	"github.com/nats-io/nats.go/micro"
)

// concrete wrapper for stock add request scope so methods can work with concrete request fields
type stockAddScope struct {
	*requestScope[schemas.StockAddRequest]
}

func (app *App) stockAddHandler(ctx context.Context, req micro.Request) {
	raw := newRequestScope[schemas.StockAddRequest](ctx, req, app.nc, app.db)
	rs := &stockAddScope{raw}
	defer rs.close(ctx)
	rs.validateJSON(ctx, app.compiler, req.Data(), schemas.StockAddRequestSchema)
	rs.decodeRequest(ctx)
	resp := rs.makeStockAddResponse(ctx, rs.addStock(ctx))
	rs.commitOrRollback(ctx)
	rs.respondJSON(ctx, req, resp)
}

// addStock adds stock to the inventory.
func (rs *stockAddScope) addStock(ctx context.Context) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "add stock")
	defer span.End()

	if rs.hasError() {
		return nil
	}

	// get typed request
	request := rs.Request()
	// if decode failed, hasError would have been true earlier and we'd have returned
	params := db.AddInventoryParams{
		ProductSku: string(request.ProductSKU),
		StockLevel: int32(request.Quantity),
	}
	inventory, err := rs.queries.AddInventory(ctx, params)
	if err != nil {
		rs.addCallerError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *stockAddScope) makeStockAddResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockAddResponse {
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
