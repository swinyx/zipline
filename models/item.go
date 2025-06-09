package models

// Item represents an item in an order with its product ID and quantity.
type Item struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
