package product

// Product mirrors the `products` table in the Northwind schema, minus the
// `supplier_ids` and `attachments` columns which are not useful over a JSON API.
type Product struct {
	ID                     int64   `json:"id"`
	ProductCode            *string `json:"product_code,omitempty"`
	ProductName            *string `json:"product_name,omitempty"`
	Description            *string `json:"description,omitempty"`
	StandardCost           float64 `json:"standard_cost"`
	ListPrice              float64 `json:"list_price"`
	ReorderLevel           *int64  `json:"reorder_level,omitempty"`
	TargetLevel            *int64  `json:"target_level,omitempty"`
	QuantityPerUnit        *string `json:"quantity_per_unit,omitempty"`
	Discontinued           bool    `json:"discontinued"`
	MinimumReorderQuantity *int64  `json:"minimum_reorder_quantity,omitempty"`
	Category               *string `json:"category,omitempty"`
}
