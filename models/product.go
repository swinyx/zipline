package models

// Product represents a product in the inventory with its ID, name, and mass in grams.
type Product struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	MassG       int    `json:"mass_g"`
}
