package store

import "zipline/models"

// Store interfaces define the methods for managing products, inventory, and orders in a store system.
type ProductStore interface {
	InitCatalog([]models.Product)
	GetProduct(int) (models.Product, bool)
	GetProductList() []models.Product
}

// InventoryStore defines methods for managing product stock in the inventory.
type InventoryStore interface {
	AddStock(productID, quantity int)
	DeductStock(productID, quantity int)
	GetStock(productID int) int
}

// OrderStore defines methods for managing pending orders in the store.
type OrderStore interface {
	AddPendingOrder(order models.Order)
	GetPendingOrders() []models.Order
	ClearPendingOrders()
}
