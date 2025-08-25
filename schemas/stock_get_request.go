package schemas

const (
	StockGetRequestSchema = "http://github.com/davidoram/beaker/schemas/stock-get.request.json"
)

// StockGetRequest represents the request structure for getting stock information.
// It corresponds to the stock-get.request.json schema.
type StockGetRequest struct {
	ProductSKU string `json:"product-sku" validate:"required"`
}
