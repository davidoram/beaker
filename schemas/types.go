package schemas

// ProductSKU represents a Stock Keeping Unit identifier for a product.
// It corresponds to the product-sku.json schema.
type ProductSKU string

// String returns the string representation of the ProductSKU.
func (sku ProductSKU) String() string {
	return string(sku)
}
