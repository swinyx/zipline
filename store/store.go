package store

import "zipline/models"

// MemoryStore groups Inventory, Catalog, and Pending Orders
type MemoryStore struct {
	// --- Inventory Store ---
	Inventory map[int]int // product_id -> quantity

	// --- Product Store ---
	Catalog map[int]models.Product // product_id -> product info

	// --- Order Store ---
	PendingOrders []models.Order // pending (unshipped) orders
}

var _ ProductStore = (*MemoryStore)(nil)
var _ InventoryStore = (*MemoryStore)(nil)
var _ OrderStore = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Inventory:     make(map[int]int),
		Catalog:       make(map[int]models.Product),
		PendingOrders: make([]models.Order, 0),
	}
}

// -----------------------------
// ProductStore methods
// -----------------------------

func (s *MemoryStore) InitCatalog(products []models.Product) {
	for _, p := range products {
		s.Catalog[p.ProductID] = p
		s.Inventory[p.ProductID] = 0
	}
}

func (s *MemoryStore) GetProduct(productID int) (models.Product, bool) {
	p, ok := s.Catalog[productID]
	return p, ok
}

func (s *MemoryStore) GetProductList() []models.Product {
	products := make([]models.Product, 0, len(s.Catalog))
	for _, product := range s.Catalog {
		products = append(products, product)
	}
	return products
}

// -----------------------------
// InventoryStore methods
// -----------------------------

func (s *MemoryStore) AddStock(productID, quantity int) {
	s.Inventory[productID] += quantity
}

func (s *MemoryStore) GetStock(productID int) int {
	return s.Inventory[productID]
}

func (s *MemoryStore) DeductStock(productID, quantity int) {
	if s.Inventory[productID] < quantity {
		// Optional: panic, log, or guard
		return
	}

	s.Inventory[productID] -= quantity
}

// -----------------------------
// OrderStore methods
// -----------------------------

func (s *MemoryStore) AddPendingOrder(order models.Order) {
	s.PendingOrders = append(s.PendingOrders, order)
}

func (s *MemoryStore) GetPendingOrders() []models.Order {
	return s.PendingOrders
}

func (s *MemoryStore) ClearPendingOrders() {
	s.PendingOrders = nil
}
