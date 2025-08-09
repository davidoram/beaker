package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/micro"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
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
func NewRequestScope(ctx context.Context, req micro.Request, pool *pgxpool.Pool) (*requestScope, error) {
	rs := &requestScope{
		req: req,
	}
	err := rs.setupDbConn(ctx, pool)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// setupDbConn establishes a connection through the pgxpool.Pool, and wraps it into a Queries instance
// which is then able to be used to access the database
func (rs *requestScope) setupDbConn(ctx context.Context, pool *pgxpool.Pool) error {
	_, span := otel.Tracer("").Start(ctx, "setup db conn")
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
	if rs.tx != nil {
		rs.tx.Rollback(ctx)
	}
	if rs.conn != nil {
		rs.conn.Release()
	}
}

// AddError adds an error to the request scope.
// If an error already exists, it appends the new error to the existing one.
// This allows for accumulating multiple errors that may occur during processing.
func (rs *requestScope) AddError(err error) {
	if rs.err == nil {
		rs.err = err
	} else {
		rs.err = fmt.Errorf("%w; %v", rs.err, err)
	}
}

func (rs *requestScope) HasError() bool { return rs.err != nil }

func (rs *requestScope) GetError() error { return rs.err }

// ValidateRequest checks if the request is valid.
// It checks if the request is nil and if the request method is valid.
// If the request is invalid, it adds an error to the request scope.
func (rs *requestScope) ValidateJSON(ctx context.Context, compiler *jsonschema.Compiler, jsonData []byte, schemaName string) {
	_, span := otel.Tracer("").Start(ctx, "ValidateJSON")
	defer span.End()

	if rs.err != nil {
		return
	}
	if len(jsonData) == 0 {
		span.SetAttributes(semconv.ErrorMessage("JSON data is empty"))
		rs.AddError(errors.New("JSON data is empty"))
		return
	}
	if compiler == nil {
		span.SetAttributes(semconv.ErrorMessage("JSON schema compiler is not initialized"))
		rs.AddError(errors.New("JSON schema compiler is not initialized"))
		return
	}
	schema, err := compiler.Compile(schemaName)
	if err != nil {
		err = fmt.Errorf("failed to compile schema %s: %w", schemaName, err)
		span.SetAttributes(semconv.ErrorMessage(err.Error()))
		rs.AddError(err)
		return
	}

	if schema == nil {
		err = fmt.Errorf("schema %s not found", schemaName)
		span.SetAttributes(semconv.ErrorMessage(err.Error()))
		rs.AddError(err)
		return
	}

	var data any
	data, err = jsonschema.UnmarshalJSON(bytes.NewReader(jsonData))
	if err != nil {
		err = fmt.Errorf("failed to unmarshal JSON data: %w", err)
		span.SetAttributes(semconv.ErrorMessage(err.Error()))
		rs.AddError(err)
		return
	}

	// Validate the data against the schema
	err = schema.Validate(data)
	if err != nil {
		err = fmt.Errorf("JSON data does not conform to schema %s: %w", schemaName, err)
		span.SetAttributes(semconv.ErrorMessage(err.Error()))
		rs.AddError(err)
		return
	}
}

// DecodeRequest decodes the request data into the provided generic type T.
// It returns the decoded value of type T. If an error occurs, it adds the error to the requestScope and returns the zero value of T.
func DecodeRequest[T any](ctx context.Context, rs *requestScope) T {
	_, span := otel.Tracer("").Start(ctx, "decode request")
	defer span.End()

	var decodedRequest T
	if rs.err != nil {
		return decodedRequest
	}

	err := json.Unmarshal(rs.req.Data(), &decodedRequest)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		rs.AddError(err)
		return decodedRequest
	}

	return decodedRequest
}

// CommitOrRollback commits the current database transaction if we have no errors.
// If there are errors, it rolls back the transaction.
// It should be called at the end of processing a request to ensure that the database state is consistent.
func (rs *requestScope) CommitOrRollback(ctx context.Context) {
	_, span := otel.Tracer("").Start(ctx, "commit or rollback")
	defer span.End()

	if rs.conn == nil || rs.tx == nil {
		return
	}
	if rs.err != nil {
		err := rs.tx.Rollback(ctx)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		return
	}
	rs.err = rs.tx.Commit(ctx)
	if rs.err != nil {
		span.SetStatus(codes.Error, rs.err.Error())
		return
	}
}

// AddStock adds stock to the inventory.
func (rs *requestScope) AddStock(ctx context.Context, stock schemas.StockAddRequest) *Inventory {
	_, span := otel.Tracer("").Start(ctx, "add stock")
	defer span.End()

	if rs.err != nil {
		return nil
	}

	params := AddInventoryParams{
		ProductSku: stock.ProductSKU.String(),
		StockLevel: int32(stock.Quantity),
	}
	inventory, err := rs.queries.AddInventory(ctx, params)
	if err != nil {
		rs.AddError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil
	}
	return &inventory
}

func (rs *requestScope) BuildStockAddResponse(inventory *Inventory) schemas.StockAddResponse {
	resp := schemas.StockAddResponse{}
	if rs.HasError() {
		resp.OK = false
		resp.Error = Ptr(rs.GetError().Error())
	} else {
		resp.OK = true
		resp.ProductSKU = Ptr(schemas.ProductSKU(inventory.ProductSku))
		resp.Quantity = Ptr(int(inventory.StockLevel))
	}
	return resp
}
