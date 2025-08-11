package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/micro"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.opentelemetry.io/otel/codes"
)

// requestScope represents a single request.
// It holds the request and any errors that occur during processing.
// Each function that processes a request should create a new requestScope instance
// When a function is called, it should check if rs.err is nil before proceeding
// If rs.err is not nil, it means an error has occurred and the function should return immediately
// This allows for early exit from the function without further processing
type requestScope struct {
	req micro.Request
	err error

	conn    *pgxpool.Conn
	tx      pgx.Tx
	queries *Queries
}

// NewRequestScope creates a new requestScope instance. It should be paired with a call to rs.Close(ctx) to guarantee cleanup.
func NewRequestScope(ctx context.Context, req micro.Request, pool *pgxpool.Pool) *requestScope {
	rs := &requestScope{
		req: req,
	}
	err := rs.setupDbConn(ctx, pool)
	if err != nil {
		rs.AddSystemError(ctx, err)
	}
	return rs
}

// setupDbConn establishes a connection through the pgxpool.Pool, and wraps it into a Queries instance
// which is then able to be used to access the database
func (rs *requestScope) setupDbConn(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, span := tracer.Start(ctx, "setup db conn")
	defer span.End()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	rs.conn = conn
	rs.tx, err = conn.Begin(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	rs.queries = New(rs.tx)
	return nil
}

func (rs *requestScope) Close(ctx context.Context) {
	if rs.conn == nil {
		return
	}
	defer func() { rs.conn = nil }()
	rs.CommitOrRollback(ctx)
	if rs.conn != nil {
		rs.conn.Release()
	}
}

// AddCallerError adds a 'caller' error to the request scope. Caller errors represent problems made by
// the API caller. It only stores the first error encountered, but it logs all errors
func (rs *requestScope) AddCallerError(ctx context.Context, err error) {
	rs.addError(ctx, err, false)
}

// AddSystemError adds a 'system' error to the request scope. System errors represent problems that occur within the system
// It only stores the first error encountered, but it logs all errors
func (rs *requestScope) AddSystemError(ctx context.Context, err error) {
	rs.addError(ctx, err, true)
}

func (rs *requestScope) addError(ctx context.Context, err error, isSystemError bool) {
	ctx, span := tracer.Start(ctx, "add error")
	defer span.End()

	if isSystemError {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "system error occurred", "error", err)
	} else {
		slog.InfoContext(ctx, "caller error occurred", "error", err)
	}

	if rs.err == nil {
		rs.err = err
	}
}

func (rs *requestScope) HasError() bool { return rs.err != nil }

func (rs *requestScope) GetError() error { return rs.err }

// ValidateRequest checks if the request is valid.
// It checks if the request is nil and if the request method is valid.
// If the request is invalid, it adds an error to the request scope.
func (rs *requestScope) ValidateJSON(ctx context.Context, compiler *jsonschema.Compiler, jsonData []byte, schemaName string) {
	ctx, span := tracer.Start(ctx, "validate JSON")
	defer span.End()

	if rs.err != nil {
		return
	}
	if len(jsonData) == 0 {
		rs.AddCallerError(ctx, errors.New("JSON data is empty"))
		return
	}
	if compiler == nil {
		rs.AddSystemError(ctx, errors.New("JSON schema compiler is not initialized"))
		return
	}
	schema, err := compiler.Compile(schemaName)
	if err != nil {
		rs.AddSystemError(ctx, fmt.Errorf("failed to compile schema %s: %w", schemaName, err))
		return
	}

	if schema == nil {
		err = fmt.Errorf("schema %s not found", schemaName)
		rs.AddSystemError(ctx, err)
		return
	}

	var data any
	data, err = jsonschema.UnmarshalJSON(bytes.NewReader(jsonData))
	if err != nil {
		rs.AddCallerError(ctx, fmt.Errorf("failed to unmarshal JSON data: %w", err))
		return
	}

	// Validate the data against the schema
	err = schema.Validate(data)
	if err != nil {
		rs.AddCallerError(ctx, fmt.Errorf("JSON data does not conform to schema %s: %w", schemaName, err))
		return
	}
}

// DecodeRequest decodes the request data into the provided generic type T.
// It returns the decoded value of type T. If an error occurs, it adds the error to the requestScope and returns the zero value of T.
func DecodeRequest[T any](ctx context.Context, rs *requestScope) T {
	ctx, span := tracer.Start(ctx, "decode request")
	defer span.End()

	var decodedRequest T
	if rs.err != nil {
		return decodedRequest
	}

	err := json.Unmarshal(rs.req.Data(), &decodedRequest)
	if err != nil {
		rs.AddCallerError(ctx, err)
		return decodedRequest
	}

	return decodedRequest
}

// CommitOrRollback commits the current database transaction if we have no errors.
// If there are errors, it rolls back the transaction.
// It should be called just before a response is sent back to the caller, so we have a chance to notify them if an error occured
func (rs *requestScope) CommitOrRollback(ctx context.Context) {

	// No transaction -> nothing to commit or rollback
	if rs.tx == nil {
		return
	}

	msg := "commit"
	if rs.HasError() {
		msg = "rollback"
	}
	ctx, span := tracer.Start(ctx, msg)
	defer span.End()

	// The transaction is unavailable after calling this function
	defer func() { rs.tx = nil }()

	// If we encountered an error during the request we need to roll back the transaction
	if rs.HasError() {
		err := rs.tx.Rollback(ctx)
		if err != nil {
			rs.AddSystemError(ctx, err)
		}
		return
	}

	// No errors, so commit the transaction
	err := rs.tx.Commit(ctx)
	if err != nil {
		rs.AddSystemError(ctx, err)
	}
}

func (rs *requestScope) RespondJSON(ctx context.Context, req micro.Request, response schemas.APIResponse) error {
	_, span := tracer.Start(ctx, "respond JSON")
	defer span.End()
	if rs.HasError() {
		response.SetErrorAttributes(rs.GetError())
	}
	return req.RespondJSON(response)
}

// AddStock adds stock to the inventory.
func (rs *requestScope) AddStock(ctx context.Context, req schemas.StockAddRequest) *Inventory {
	ctx, span := tracer.Start(ctx, "add stock")
	defer span.End()

	if rs.err != nil {
		return nil
	}

	params := AddInventoryParams{
		ProductSku: req.ProductSKU.String(),
		StockLevel: int32(req.Quantity),
	}
	inventory, err := rs.queries.AddInventory(ctx, params)
	if err != nil {
		rs.AddCallerError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *requestScope) MakeStockAddResponse(ctx context.Context, inventory *Inventory) *schemas.StockAddResponse {
	ctx, span := tracer.Start(ctx, "build stock-add response")
	defer span.End()

	resp := schemas.StockAddResponse{}
	if rs.HasError() {
		resp.OK = false
		resp.Error = Ptr(rs.GetError().Error())
	} else {
		resp.OK = true
		resp.ProductSKU = Ptr(schemas.ProductSKU(inventory.ProductSku))
		resp.Quantity = Ptr(int(inventory.StockLevel))
	}
	return &resp
}

// GetStock retrieves the stock information for a product.
func (rs *requestScope) GetStock(ctx context.Context, req schemas.StockGetRequest) *Inventory {
	ctx, span := tracer.Start(ctx, "get stock")
	defer span.End()

	if rs.err != nil {
		return nil
	}

	inventory, err := rs.queries.GetInventory(ctx, req.ProductSKU.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "no inventory found for product", "product_sku", req.ProductSKU.String())
			return &Inventory{ProductSku: req.ProductSKU.String(), StockLevel: 0}
		}
		rs.AddSystemError(ctx, err)
		return nil
	}
	return &inventory
}

func (rs *requestScope) MakeStockGetResponse(ctx context.Context, inventory *Inventory) *schemas.StockGetResponse {
	ctx, span := tracer.Start(ctx, "build stock-get response")
	defer span.End()

	resp := schemas.StockGetResponse{}
	if rs.HasError() {
		resp.OK = false
		resp.Error = Ptr(rs.GetError().Error())
	} else {
		resp.OK = true
		resp.ProductSKU = Ptr(schemas.ProductSKU(inventory.ProductSku))
		resp.Quantity = Ptr(int(inventory.StockLevel))
	}
	return &resp
}

// RemoveStock adds stock to the inventory.
func (rs *requestScope) RemoveStock(ctx context.Context, req schemas.StockRemoveRequest) *Inventory {
	ctx, span := tracer.Start(ctx, "remove stock")
	defer span.End()

	if rs.err != nil {
		return nil
	}

	params := RemoveInventoryParams{
		ProductSku: req.ProductSKU.String(),
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

func (rs *requestScope) MakeStockRemoveResponse(ctx context.Context, inventory *Inventory) *schemas.StockRemoveResponse {
	_, span := tracer.Start(ctx, "build stock-remove response")
	defer span.End()

	resp := schemas.StockRemoveResponse{}
	if rs.HasError() {
		resp.OK = false
		resp.Error = Ptr(rs.GetError().Error())
	} else {
		resp.OK = true
		resp.ProductSKU = Ptr(schemas.ProductSKU(inventory.ProductSku))
		resp.Quantity = Ptr(int(inventory.StockLevel))
	}
	return &resp
}
