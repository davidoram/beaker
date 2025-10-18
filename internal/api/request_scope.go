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
type requestScope[T any] struct {
	nc  *nats.Conn
	req micro.Request
	err error

	conn    *pgxpool.Conn
	tx      pgx.Tx
	queries *db.Queries

	// decoded holds the decoded request body as the concrete type T.
	decoded T
}

// newRequestScope creates a new requestScope instance. It should be paired with a call to rs.Close(ctx) to guarantee cleanup.

func newRequestScope[T any](ctx context.Context, req micro.Request, nc *nats.Conn, pool *pgxpool.Pool) *requestScope[T] {
	rs := &requestScope[T]{
		req: req,
		nc:  nc,
	}
	rs.setupDbConn(ctx, pool)
	return rs
}

func (rs *requestScope[T]) close(ctx context.Context) {
	if rs.conn == nil {
		return
	}
	defer func() { rs.conn = nil }()
	rs.commitOrRollback(ctx)
	if rs.conn != nil {
		rs.conn.Release()
	}
}

// setupDbConn establishes a connection through the pgxpool.Pool, and wraps it into a Queries instance
// which is then able to be used to access the database
func (rs *requestScope[T]) setupDbConn(ctx context.Context, pool *pgxpool.Pool) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "setup db conn")
	defer span.End()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		rs.addSystemError(ctx, err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	rs.conn = conn
	rs.tx, err = conn.Begin(ctx)
	if err != nil {
		rs.addSystemError(ctx, err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	rs.queries = db.New(rs.tx)
}

// addCallerError adds a 'caller' error to the request scope which represents a problem made by
// the API caller.
func (rs *requestScope[T]) addCallerError(ctx context.Context, err error) {
	rs.addError(ctx, err, false)
}

// addSystemError adds a 'system' error to the request scope which represents a problem that occurs inside our system
func (rs *requestScope[T]) addSystemError(ctx context.Context, err error) {
	rs.addError(ctx, err, true)
}

// adds an error, it only stores the first error encountered, but it logs all errors.
// If its a system error will log at error level, and mark the span in error because as system owners we need
// to be aware of these errors. Caller errors are logged at info level and do not mark the span in error
func (rs *requestScope[T]) addError(ctx context.Context, err error, isSystemError bool) {
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

func (rs *requestScope[T]) hasError() bool { return rs.err != nil }

func (rs *requestScope[T]) getError() error { return rs.err }

// ValidateRequest checks if the request is valid.
// It checks if the request is nil and if the request method is valid.
// If the request is invalid, it adds an error to the request scope.
func (rs *requestScope[T]) validateJSON(ctx context.Context, compiler *jsonschema.Compiler, jsonData []byte, schemaName string) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "validate JSON")
	defer span.End()

	if rs.hasError() {
		return
	}
	if len(jsonData) == 0 {
		rs.addCallerError(ctx, errors.New("JSON data is empty"))
		return
	}
	if compiler == nil {
		rs.addSystemError(ctx, errors.New("JSON schema compiler is not initialized"))
		return
	}
	schema, err := compiler.Compile(schemaName)
	if err != nil {
		rs.addSystemError(ctx, fmt.Errorf("failed to compile schema %s: %w", schemaName, err))
		return
	}

	if schema == nil {
		rs.addSystemError(ctx, fmt.Errorf("schema %s not found", schemaName))
		return
	}

	var data any
	data, err = jsonschema.UnmarshalJSON(bytes.NewReader(jsonData))
	if err != nil {
		rs.addCallerError(ctx, fmt.Errorf("failed to unmarshal JSON data: %w", err))
		return
	}

	// Validate the data against the schema
	err = schema.Validate(data)
	if err != nil {
		rs.addCallerError(ctx, fmt.Errorf("JSON data does not conform to schema %s: %w", schemaName, err))
		return
	}
}

// decodeRequest unmarshals the request data into the decoded value. If an error occurs, it adds the error to the requestScope.
func (rs *requestScope[T]) decodeRequest(ctx context.Context) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "decode request")
	defer span.End()

	if rs.hasError() {
		return
	}

	var decodedRequest T
	err := json.Unmarshal(rs.req.Data(), &decodedRequest)
	if err != nil {
		rs.addCallerError(ctx, err)
	}

	// store decoded request on the scope for later use
	rs.decoded = decodedRequest
}

// Request returns the decoded request value (zero value if decode failed).
// Use rs.hasError() to determine whether decoding reported an error.
func (rs *requestScope[T]) Request() T {
	return rs.decoded
}

// commitOrRollback commits the current database transaction if we have no errors.
// If there are errors, it rolls back the transaction.
// It should be called just before a response is sent back to the caller, so we have a chance to notify them if an error occured
func (rs *requestScope[T]) commitOrRollback(ctx context.Context) {

	// No transaction -> nothing to commit or rollback
	if rs.tx == nil {
		return
	}

	msg := "tx commit"
	if rs.hasError() {
		msg = "tx rollback"
	}
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, msg)
	defer span.End()

	// The transaction is unavailable after calling this function
	defer func() { rs.tx = nil }()

	// If we encountered an error during the request we need to roll back the transaction
	if rs.hasError() {
		err := rs.tx.Rollback(ctx)
		if err != nil {
			rs.addSystemError(ctx, err)
		}
		return
	}

	// No errors, so commit the transaction
	err := rs.tx.Commit(ctx)
	if err != nil {
		rs.addSystemError(ctx, err)
	}
}

func (rs *requestScope[T]) respondJSON(ctx context.Context, req micro.Request, response schemas.APIResponse) {
	tracer := telemetry.GetTracer()
	_, span := tracer.Start(ctx, "respond JSON")
	defer span.End()
	if rs.hasError() {
		slog.ErrorContext(ctx, "Request has error", "error", rs.getError())
		response.SetErrorAttributes(rs.getError())
	}
	err := req.RespondJSON(response)
	if err != nil {
		response.SetErrorAttributes(rs.getError())
		slog.ErrorContext(ctx, "RespondJSON returned error", "error", err)
	}
}

func (rs *requestScope[T]) emitEvent(ctx context.Context, event schemas.LowStockEvent) error {
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
	return rs.nc.Publish(event.Subject(), []byte(eventJSON))
}

// emitLowStockEvent checks if the updated inventory is below the low stock threshold
func (rs *requestScope[T]) emitLowStockEvent(ctx context.Context, updatedInventory *db.Inventory) {

	if rs.hasError() {
		return
	}
	// If stock was successfully removed and is now low, emit a LowStockEvent
	if updatedInventory.StockLevel < LowStockThreshold {
		event := schemas.LowStockEvent{
			ProductSKU: updatedInventory.ProductSku,
			StockLevel: int(updatedInventory.StockLevel),
		}
		if err := rs.emitEvent(ctx, event); err != nil {
			rs.addSystemError(ctx, err)
		}
	}
}
