package schemas

// Helper functions for creating response structures

// NewStockAddSuccessResponse creates a successful StockAddResponse.
func NewStockAddSuccessResponse(productSKU string, quantity int) StockAddResponse {
	return StockAddResponse{
		ProductSKU: &productSKU,
		Quantity:   &quantity,
	}
}

// NewStockAddErrorResponse creates an error StockAddResponse.
func NewStockAddErrorResponse(errorMsg string) StockAddResponse {
	return StockAddResponse{
		Error: &errorMsg,
	}
}

// NewStockRemoveSuccessResponse creates a successful StockRemoveResponse.
func NewStockRemoveSuccessResponse(productSKU string, quantity int) StockRemoveResponse {
	return StockRemoveResponse{
		ProductSKU: &productSKU,
		Quantity:   &quantity,
	}
}

// NewStockRemoveErrorResponse creates an error StockRemoveResponse.
func NewStockRemoveErrorResponse(errorMsg string) StockRemoveResponse {
	return StockRemoveResponse{
		Error: &errorMsg,
	}
}

// NewStockGetSuccessResponse creates a successful StockGetResponse.
func NewStockGetSuccessResponse(productSKU string, quantity int) StockGetResponse {
	return StockGetResponse{
		ProductSKU: &productSKU,
		Quantity:   &quantity,
	}
}

// NewStockGetErrorResponse creates an error StockGetResponse.
func NewStockGetErrorResponse(errorMsg string) StockGetResponse {
	return StockGetResponse{
		Error: &errorMsg,
	}
}
