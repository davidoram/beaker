package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/davidoram/beaker/schemas"
	"github.com/nats-io/nats.go/micro"
)

// fakeRequest implements the micro.Request interface for tests
type fakeRequest struct {
	data []byte
}

func (f *fakeRequest) Respond(response []byte, _ ...micro.RespondOpt) error { return nil }
func (f *fakeRequest) RespondJSON(v any, _ ...micro.RespondOpt) error       { return nil }
func (f *fakeRequest) Error(code, description string, _ []byte, _ ...micro.RespondOpt) error {
	return nil
}
func (f *fakeRequest) Data() []byte           { return f.data }
func (f *fakeRequest) Headers() micro.Headers { return nil }
func (f *fakeRequest) Subject() string        { return "test.subject" }
func (f *fakeRequest) Reply() string          { return "" }

// Test happy path decoding into StockGetRequest
func TestDecodeRequest_HappyPath(t *testing.T) {
	ctx := context.Background()
	reqObj := schemas.StockGetRequest{ProductSKU: "SKU123"}
	b, err := json.Marshal(reqObj)
	if err != nil {
		t.Fatalf("failed to marshal test request: %v", err)
	}

	req := &fakeRequest{data: b}
	// instantiate requestScope directly to avoid DB pool setup
	rs := &requestScope[schemas.StockGetRequest]{req: req}
	defer rs.close(ctx)

	_ = rs.decodeRequest(ctx)
	if rs.hasError() {
		t.Fatalf("expected no error decoding valid JSON, got: %v", rs.getError())
	}

	decoded := rs.Request()
	if decoded.ProductSKU != "SKU123" {
		t.Fatalf("decoded value mismatch: got %v", decoded)
	}
}

// Test invalid JSON sets an error on the scope
func TestDecodeRequest_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	b := []byte("not a json")
	req := &fakeRequest{data: b}
	rs := &requestScope[schemas.StockGetRequest]{req: req}
	defer rs.close(ctx)

	_ = rs.decodeRequest(ctx)
	if !rs.hasError() {
		t.Fatalf("expected hasError after invalid JSON decode")
	}
}
