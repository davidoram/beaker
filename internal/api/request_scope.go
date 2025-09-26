package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/davidoram/beaker/internal/db"
	"github.com/davidoram/beaker/internal/telemetry"
	"github.com/davidoram/beaker/schemas"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.opentelemetry.io/otel/codes"
)

const LowStockThreshold = 10

// requestScope holds the context for a single request.
// It holds the request and any errors that occur during processing.
// When the API receives a call it should create a NewRequestScope instance
// Functions that work through the various phases of the request, check if any preceding
// phases encountered an error by checking if rs.err is nil before proceeding.
// If rs.err is not nil, it means an error has occurred and the function should
// act appropriately.
// This allows for early exit from the function without further processing
type requestScope struct {
	nc  *nats.Conn
	req micro.Request
	err error

	conn    *pgxpool.Conn
	tx      pgx.Tx
	queries *db.Queries
}

// NewRequestScope creates a new requestScope instance. It should be paired with a call to rs.Close(ctx) to guarantee cleanup.
func NewRequestScope(ctx context.Context, req micro.Request, nc *nats.Conn, pool *pgxpool.Pool) *requestScope {
	rs := &requestScope{
		req: req,
		nc:  nc,
	}
	rs.setupDbConn(ctx, pool)
	return rs
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

// setupDbConn establishes a connection through the pgxpool.Pool, and wraps it into a Queries instance
// which is then able to be used to access the database
func (rs *requestScope) setupDbConn(ctx context.Context, pool *pgxpool.Pool) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "setup db conn")
	defer span.End()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		rs.AddSystemError(ctx, err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	rs.conn = conn
	rs.tx, err = conn.Begin(ctx)
	if err != nil {
		rs.AddSystemError(ctx, err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	rs.queries = db.New(rs.tx)
}

// AddCallerError adds a 'caller' error to the request scope which represents a problem made by
// the API caller.
func (rs *requestScope) AddCallerError(ctx context.Context, err error) {
	rs.addError(ctx, err, false)
}

// AddSystemError adds a 'system' error to the request scope which represents a problem that occurs inside our system
func (rs *requestScope) AddSystemError(ctx context.Context, err error) {
	rs.addError(ctx, err, true)
}

// adds an error, it only stores the first error encountered, but it logs all errors.
// If its a system error will log at error level, and mark the span in error because as system owners we need
// to be aware of these errors. Caller errors are logged at info level and do not mark the span in error
func (rs *requestScope) addError(ctx context.Context, err error, isSystemError bool) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "add error")
	defer span.End()

	// Mark system errors so that they will can be filtered easily inside OpenTelemetry
	if isSystemError {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "system error", "error", err)
	} else {
		slog.InfoContext(ctx, "caller error", "error", err)
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
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "validate JSON")
	defer span.End()

	if rs.HasError() {
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
		rs.AddSystemError(ctx, fmt.Errorf("schema %s not found", schemaName))
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
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "decode request")
	defer span.End()

	var decodedRequest T
	if rs.HasError() {
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

	msg := "tx commit"
	if rs.HasError() {
		msg = "tx rollback"
	}
	tracer := telemetry.GetTracer()
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

func (rs *requestScope) RespondJSON(ctx context.Context, req micro.Request, response schemas.APIResponse) {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "respond JSON")
	defer span.End()
	if rs.HasError() {
		slog.ErrorContext(ctx, "Request has error", "error", rs.GetError())
		response.SetErrorAttributes(rs.GetError())
	}
	err := req.RespondJSON(response)
	if err != nil {
		response.SetErrorAttributes(rs.GetError())
		slog.ErrorContext(ctx, "RespondJSON returned error", "error", err)
	}
}

func (rs *requestScope) EmitEvent(ctx context.Context, event schemas.LowStockEvent) error {
	log.Printf("Emitting low stock event: %+v", event)
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "emit low stock event")
	defer span.End()
	slog.InfoContext(ctx, "Emitting low stock event", "event", event)
	eventJSON, err := json.Marshal(event)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "Failed to marshal low stock event", "error", err)
		return err
	}
	rs.nc.Publish(event.Subject(), []byte(eventJSON))
	return nil
}

// EmitLowStockEvent checks if the updated inventory is below the low stock threshold
func (rs *requestScope) EmitLowStockEvent(ctx context.Context, updatedInventory *db.Inventory) {

	if rs.HasError() {
		return
	}
	// If stock was successfully removed and is now low, emit a LowStockEvent
	if updatedInventory.StockLevel < LowStockThreshold {
		event := schemas.LowStockEvent{
			ProductSKU: updatedInventory.ProductSku,
			StockLevel: int(updatedInventory.StockLevel),
		}
		if err := rs.EmitEvent(ctx, event); err != nil {
			rs.AddSystemError(ctx, err)
		}
	}
}
