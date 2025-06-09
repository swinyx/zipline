package store

import (
	"testing"
	"zipline/models"
)

func TestMemoryStore_InitCatalog(t *testing.T) {
	store := NewMemoryStore()
	products := []models.Product{
		{ProductID: 1, ProductName: "Product A", MassG: 500},
		{ProductID: 2, ProductName: "Product B", MassG: 300},
	}
	store.InitCatalog(products)

	if len(store.Catalog) != 2 {
		t.Errorf("expected catalog length 2, got %d", len(store.Catalog))
	}
	if store.Inventory[1] != 0 || store.Inventory[2] != 0 {
		t.Errorf("expected initial inventory to be 0 for all products")
	}
}

func TestMemoryStore_GetProduct(t *testing.T) {
	store := NewMemoryStore()
	product := models.Product{ProductID: 1, ProductName: "Test", MassG: 100}
	store.InitCatalog([]models.Product{product})

	p, found := store.GetProduct(1)
	if !found {
		t.Fatal("expected product to be found")
	}
	if p.ProductName != "Test" {
		t.Errorf("expected 'Test', got %s", p.ProductName)
	}
}

func TestMemoryStore_GetProductList(t *testing.T) {
	store := NewMemoryStore()
	store.InitCatalog([]models.Product{
		{ProductID: 1, ProductName: "Product X", MassG: 700},
		{ProductID: 2, ProductName: "Product Y", MassG: 800},
	})

	list := store.GetProductList()
	if len(list) != 2 {
		t.Errorf("expected 2 products, got %d", len(list))
	}
}

func TestMemoryStore_InventoryOperations(t *testing.T) {
	store := NewMemoryStore()
	store.AddStock(1, 10)
	if stock := store.GetStock(1); stock != 10 {
		t.Errorf("expected stock 10, got %d", stock)
	}

	store.DeductStock(1, 3)
	if stock := store.GetStock(1); stock != 7 {
		t.Errorf("expected stock 7, got %d", stock)
	}
}

func TestMemoryStore_OrderOperations(t *testing.T) {
	store := NewMemoryStore()
	order := models.Order{
		OrderID: 101,
		Requested: []models.Item{
			{ProductID: 1, Quantity: 2},
		},
	}

	store.AddPendingOrder(order)

	orders := store.GetPendingOrders()
	if len(orders) != 1 {
		t.Errorf("expected 1 pending order, got %d", len(orders))
	}
	if orders[0].OrderID != 101 {
		t.Errorf("expected order ID 101, got %d", orders[0].OrderID)
	}

	store.ClearPendingOrders()
	if len(store.GetPendingOrders()) != 0 {
		t.Error("expected 0 pending orders after clearing")
	}
}
