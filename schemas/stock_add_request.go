package schemas

const (
	StockAddRequestSchema = "http://github.com/davidoram/beaker/schemas/stock-add.request.json"
)

// StockAddRequest represents the request structure for adding stock.
// It corresponds to the stock-add.request.json schema.
type StockAddRequest struct {
	ProductSKU ProductSKU `json:"product-sku"`
	Quantity   int        `json:"quantity"`
}
