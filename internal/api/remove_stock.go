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
	raw := newRequestScope[schemas.StockRemoveRequest](ctx, req, app.nc, app.db)
	rs := &stockRemoveScope{raw}
	defer rs.close(ctx)
	rs.validateJSON(ctx, app.compiler, req.Data(), schemas.StockRemoveRequestSchema)
	_ = rs.decodeRequest(ctx)
	if rs.hasError() {
		resp := rs.makeStockRemoveResponse(ctx, nil)
		rs.commitOrRollback(ctx)
		rs.respondJSON(ctx, req, resp)
		return
	}
	updatedInventory := rs.removeStock(ctx)
	rs.emitLowStockEvent(ctx, updatedInventory)
	resp := rs.makeStockRemoveResponse(ctx, updatedInventory)
	rs.commitOrRollback(ctx)
	rs.respondJSON(ctx, req, resp)
}

// removeStock adds stock to the inventory.
// concrete wrapper for stock remove request scope
type stockRemoveScope struct {
	*requestScope[schemas.StockRemoveRequest]
}

func (rs *stockRemoveScope) removeStock(ctx context.Context) *db.Inventory {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "remove stock")
	defer span.End()

	if rs.hasError() {
		return nil
	}

	reqTyped := rs.Request()
	params := db.RemoveInventoryParams{
		ProductSku: reqTyped.ProductSKU,
		StockLevel: int32(reqTyped.Quantity),
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
					rs.addCallerError(ctx, fmt.Errorf("stock level cannot go below zero for %s", reqTyped.ProductSKU))
				case "inventory_product_sku_format":
					rs.addCallerError(ctx, fmt.Errorf("invalid SKU format: %s", reqTyped.ProductSKU))
				default:
					rs.addCallerError(ctx, fmt.Errorf("business rule violated: %s", pgErr.Message))
				}
				return nil
			}
		}
		rs.addSystemError(ctx, fmt.Errorf("database error: %s", err.Error()))
		return nil
	}
	return &inventory
}

func (rs *stockRemoveScope) makeStockRemoveResponse(ctx context.Context, inventory *db.Inventory) *schemas.StockRemoveResponse {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "build stock-remove response")
	defer span.End()

	resp := schemas.StockRemoveResponse{}
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
