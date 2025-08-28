package schemas

import "github.com/davidoram/beaker/internal/utility"

// StockGetResponse represents the response structure for getting stock information.
// It corresponds to the stock-get.response.json schema.
// This implements the oneOf pattern using interface{} - you should check the actual type at runtime.
type StockGetResponse struct {
	// OK is true with a successful response, false with an error response
	OK bool `json:"ok"`

	// Success response fields
	ProductSKU *string `json:"product-sku,omitempty"`
	Quantity   *int    `json:"quantity,omitempty"`

	// Error response field
	Error *string `json:"error,omitempty"`
}

func (r *StockGetResponse) SetErrorAttributes(err error) {
	r.Error = utility.Ptr(err.Error())
	r.OK = false

	r.ProductSKU = nil
	r.Quantity = nil
}
