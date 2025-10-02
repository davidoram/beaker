package schemas

const (
	LockStockEventSchema = "http://github.com/davidoram/beaker/schemas/low-stock.event.json"
)

// LowStockEvent represents the event generated when stock is low.
// It corresponds to the low-stock.event.json schema.
type LowStockEvent struct {
	ProductSKU string `json:"product-sku"`
	StockLevel int    `json:"stock-level"`
}

// Subject returns the NATS subject that LowStockEvent will be published to.
func (e LowStockEvent) Subject() string {
	return "events.low_stock"
}
