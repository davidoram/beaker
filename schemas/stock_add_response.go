package schemas

import "github.com/davidoram/beaker/internal/utility"

// StockAddResponse represents the response structure for adding stock.
// It corresponds to the stock-add.response.json schema.
// This implements the oneOf pattern using interface{} - you should check the actual type at runtime.
type StockAddResponse struct {
	// OK is true with a successful response, false with an error response
	OK bool `json:"ok"`

	// Success response fields
	ProductSKU *string `json:"product-sku,omitempty"`
	Quantity   *int    `json:"quantity,omitempty"`

	// Error response field
	Error *string `json:"error,omitempty"`
}

func (r *StockAddResponse) SetErrorAttributes(err error) {
	r.Error = utility.Ptr(err.Error())
	r.OK = false

	r.ProductSKU = nil
	r.Quantity = nil
}
