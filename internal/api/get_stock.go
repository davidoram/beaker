package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/davidoram/beaker/internal/db"
	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/internal/utility"
	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go/micro"
)

func (app *App) stockGetHandler(ctx context.Context, req micro.Request) {
	rs := newRequestScope(ctx, req, app.nc, app.db)
	defer rs.close(ctx)
	rs.validateJSON(ctx, app.compiler, req.Data(), schemas.StockGetRequestSchema)
	stockReq := decodeRequest[schemas.StockGetRequest](ctx, rs)
	resp := rs.makeStockGetResponse(ctx, rs.getStock(ctx, stockReq))
	rs.commitOrRollback(ctx)
	rs.respondJSON(ctx, req, resp)
}

// getStock retrieves the stock information for a product.
func (rs *requestScope) getStock(ctx context.Context, req schemas.StockGetRequest) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "get stock")
	defer span.End()

	if rs.hasError() {
		return nil
	}

	inventory, err := rs.queries.GetInventory(ctx, req.ProductSKU)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "no inventory found for product", "product_sku", req.ProductSKU)
			return &db.Inventory{ProductSku: req.ProductSKU, StockLevel: 0}
		}
		rs.addSystemError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *requestScope) makeStockGetResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockGetResponse {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "build stock-get response")
	defer span.End()

	resp := schemas.StockGetResponse{}
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
