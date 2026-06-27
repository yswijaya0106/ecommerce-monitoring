package order

import "time"

// OrderSummary is a lightweight row used for listing orders.
type OrderSummary struct {
	ID         int64      `json:"id"`
	CustomerID *int64     `json:"customer_id,omitempty"`
	OrderDate  *time.Time `json:"order_date,omitempty"`
	StatusID   *int64     `json:"status_id,omitempty"`
	Notes      *string    `json:"notes,omitempty"`
	Total      float64    `json:"total"`
}

// Item is one product line within an order, joined with the product name
// so the API consumer doesn't need a second round trip.
type Item struct {
	ID          int64   `json:"id"`
	ProductID   *int64  `json:"product_id,omitempty"`
	ProductName *string `json:"product_name,omitempty"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Discount    float64 `json:"discount"`
}

// Order is the full order with its line items, returned by Get/Create.
type Order struct {
	OrderSummary
	Items []Item `json:"items"`
}

// CreateItem is one requested line item when placing an order. UnitPrice is
// resolved server-side from the product's current list_price, not trusted
// from the client.
type CreateItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

// CreateRequest is the payload for placing a new order.
type CreateRequest struct {
	CustomerID int64        `json:"customer_id"`
	Notes      *string      `json:"notes,omitempty"`
	Items      []CreateItem `json:"items"`
}
