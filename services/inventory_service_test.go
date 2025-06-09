package services_test

import (
	"testing"

	"zipline/models"
	"zipline/services"
	"zipline/store"
)

func TestShipPackage_ProductNotFound(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	// No catalog initialized
	order := models.Order{
		OrderID: 404,
		Requested: []models.Item{
			{ProductID: 999, Quantity: 1},
		},
	}

	err := svc.ShipPackage(order)
	if err == nil {
		t.Fatalf("Expected error for unknown product, got nil")
	}
}

func TestInitCatalog_EmptyProductInfo(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	err := svc.InitCatalog([]models.Product{})
	if err == nil {
		t.Errorf("Expected error when initializing with empty product list, got nil")
	}
}

func TestProcessOrder_FullAndPartialShipment(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	// Step 1: Initialize product catalog
	catalog := []models.Product{
		{ProductID: 0, ProductName: "RBC A+ Adult", MassG: 700},
		{ProductID: 10, ProductName: "FFP A+", MassG: 300},
	}
	svc.InitCatalog(catalog)

	// Step 2: Restock some products (only enough for partial fulfillment)
	restock := []models.Item{
		{ProductID: 0, Quantity: 2},  // 2 x 700g = 1400g
		{ProductID: 10, Quantity: 1}, // 1 x 300g = 300g
	}
	svc.ProcessRestock(restock)

	// Step 3: Submit an order that exceeds current stock
	order := models.Order{
		OrderID: 123,
		Requested: []models.Item{
			{ProductID: 0, Quantity: 2},  // fully in stock
			{ProductID: 10, Quantity: 3}, // only partially in stock
		},
	}
	svc.ProcessOrder(order)

	// Step 4: Check that the pending order was correctly saved
	pending := memStore.GetPendingOrders()
	if len(pending) != 1 {
		t.Fatalf("Expected 1 pending order, got %d", len(pending))
	}

	if pending[0].OrderID != 123 {
		t.Errorf("Expected OrderID 123 in pending, got %d", pending[0].OrderID)
	}

	// Step 5: Restock and ensure pending order gets reprocessed
	svc.ProcessRestock([]models.Item{
		{ProductID: 10, Quantity: 3},
	})

	// After restock and reprocessing, there should be no pending orders
	if len(memStore.GetPendingOrders()) != 0 {
		t.Errorf("Expected no pending orders after restock and reprocessing")
	}
}

func TestProcessOrder_ShipmentOverweightSplit(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	svc.InitCatalog([]models.Product{
		{ProductID: 1, ProductName: "Heavy Item", MassG: 1000},
	})

	memStore.AddStock(1, 3)

	order := models.Order{
		OrderID: 1001,
		Requested: []models.Item{
			{ProductID: 1, Quantity: 3}, // 3 x 1000g = 3000g
		},
	}

	err := svc.ProcessOrder(order)
	if err != nil {
		t.Errorf("Expected successful split shipment, got error: %v", err)
	}

	if memStore.GetStock(1) != 0 {
		t.Errorf("Expected all stock to be consumed, got %d", memStore.GetStock(1))
	}
}

func TestShipPackage_ValidatesAndDeductsStock(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	catalog := []models.Product{
		{ProductID: 1, ProductName: "PLT O+", MassG: 80},
	}
	svc.InitCatalog(catalog)
	memStore.AddStock(1, 5)

	order := models.Order{
		OrderID: 555,
		Requested: []models.Item{
			{ProductID: 1, Quantity: 3},
		},
	}

	err := svc.ShipPackage(order)
	if err != nil {
		t.Fatalf("Expected shipment to succeed, got error: %v", err)
	}

	remaining := memStore.GetStock(1)
	if remaining != 2 {
		t.Errorf("Expected 2 units remaining, got %d", remaining)
	}
}

func TestShipPackage_ExceedsWeight(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	catalog := []models.Product{
		{ProductID: 2, ProductName: "FFP AB+", MassG: 1000},
	}
	svc.InitCatalog(catalog)
	memStore.AddStock(2, 2)

	order := models.Order{
		OrderID: 777,
		Requested: []models.Item{
			{ProductID: 2, Quantity: 2}, // 2000g total
		},
	}

	err := svc.ShipPackage(order)
	if err == nil {
		t.Fatalf("Expected weight limit error, but got nil")
	}
}

func TestShipPackage_InsufficientStock(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	catalog := []models.Product{
		{ProductID: 3, ProductName: "CRYO A+", MassG: 40},
	}
	svc.InitCatalog(catalog)
	memStore.AddStock(3, 1)

	order := models.Order{
		OrderID: 888,
		Requested: []models.Item{
			{ProductID: 3, Quantity: 2},
		},
	}

	err := svc.ShipPackage(order)
	if err == nil {
		t.Fatalf("Expected stock error, but got nil")
	}
}

func TestProcessRestock_NoItems(t *testing.T) {
	memStore := store.NewMemoryStore()
	svc := services.NewInventoryService(memStore)

	err := svc.ProcessRestock([]models.Item{})
	if err == nil {
		t.Errorf("Expected error when processing empty restock, got nil")
	}
}
