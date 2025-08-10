package schemas

const (
	StockRemoveRequestSchema = "http://github.com/davidoram/beaker/schemas/stock-remove.request.json"
)

// StockRemoveRequest represents the request structure for removing stock.
// It corresponds to the stock-remove.request.json schema.
type StockRemoveRequest struct {
	ProductSKU ProductSKU `json:"product-sku"`
	Quantity   int        `json:"quantity"`
}
