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
	rs := NewRequestScope(ctx, req, app.nc, app.db)
	defer rs.Close(ctx)
	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockGetRequestSchema)
	stockReq := DecodeRequest[schemas.StockGetRequest](ctx, rs)
	resp := rs.MakeStockGetResponse(ctx, rs.GetStock(ctx, stockReq))
	rs.CommitOrRollback(ctx)
	rs.RespondJSON(ctx, req, resp)
}

// GetStock retrieves the stock information for a product.
func (rs *requestScope) GetStock(ctx context.Context, req schemas.StockGetRequest) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "get stock")
	defer span.End()

	if rs.HasError() {
		return nil
	}

	inventory, err := rs.queries.GetInventory(ctx, req.ProductSKU)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "no inventory found for product", "product_sku", req.ProductSKU)
			return &db.Inventory{ProductSku: req.ProductSKU, StockLevel: 0}
		}
		rs.AddSystemError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *requestScope) MakeStockGetResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockGetResponse {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "build stock-get response")
	defer span.End()

	resp := schemas.StockGetResponse{}
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
