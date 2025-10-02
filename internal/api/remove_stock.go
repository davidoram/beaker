package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/davidoram/beaker/internal/db"
	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/internal/utility"
	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nats-io/nats.go/micro"
)

func (app *App) stockRemoveHandler(ctx context.Context, req micro.Request) {
	rs := NewRequestScope(ctx, req, app.nc, app.db)
	defer rs.Close(ctx)
	rs.ValidateJSON(ctx, app.compiler, req.Data(), schemas.StockRemoveRequestSchema)
	stockReq := DecodeRequest[schemas.StockRemoveRequest](ctx, rs)
	updatedInventory := rs.RemoveStock(ctx, stockReq)
	rs.EmitLowStockEvent(ctx, updatedInventory)
	resp := rs.MakeStockRemoveResponse(ctx, updatedInventory)
	rs.CommitOrRollback(ctx)
	rs.RespondJSON(ctx, req, resp)
}

// RemoveStock adds stock to the inventory.
func (rs *requestScope) RemoveStock(ctx context.Context, req schemas.StockRemoveRequest) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "remove stock")
	defer span.End()

	if rs.HasError() {
		return nil
	}

	params := db.RemoveInventoryParams{
		ProductSku: req.ProductSKU,
		StockLevel: int32(req.Quantity),
	}
	inventory, err := rs.queries.RemoveInventory(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Detect any CHECK violation
			if pgErr.Code == pgerrcode.CheckViolation {
				slog.InfoContext(ctx, "database constraint error", "code", pgErr.Code, "message", pgErr.Message, "constraint", pgErr.ConstraintName)
				// Branch by constraint name
				switch pgErr.ConstraintName {
				case "inventory_stock_level_nonnegative":
					rs.AddCallerError(ctx, fmt.Errorf("stock level cannot go below zero for %s", req.ProductSKU))
				case "inventory_product_sku_format":
					rs.AddCallerError(ctx, fmt.Errorf("invalid SKU format: %s", req.ProductSKU))
				default:
					rs.AddCallerError(ctx, fmt.Errorf("business rule violated: %s", pgErr.Message))
				}
				return nil
			}
		}
		rs.AddSystemError(ctx, fmt.Errorf("database error: %s", err.Error()))
		return nil
	}
	return &inventory
}

func (rs *requestScope) MakeStockRemoveResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockRemoveResponse {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "build stock-remove response")
	defer span.End()

	resp := schemas.StockRemoveResponse{}
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
