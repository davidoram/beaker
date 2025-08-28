package schemas

import "github.com/davidoram/beaker/internal/utility"

// StockRemoveResponse represents the response structure for removing stock.
// It corresponds to the stock-remove.response.json schema.
// This implements the oneOf pattern using interface{} - you should check the actual type at runtime.
type StockRemoveResponse struct {
	// OK is true with a successful response, false with an error response
	OK bool `json:"ok"`

	// Success response fields
	ProductSKU *string `json:"product-sku,omitempty"`
	Quantity   *int    `json:"quantity,omitempty"`

	// Error response field
	Error *string `json:"error,omitempty"`
}

func (r *StockRemoveResponse) SetErrorAttributes(err error) {
	r.Error = utility.Ptr(err.Error())
	r.OK = false

	r.ProductSKU = nil
	r.Quantity = nil
}
