package models

// Order represents a customer's order containing multiple items.
type Order struct {
	OrderID   int    `json:"order_id"`
	Requested []Item `json:"requested"`
}
