package main

import (
	"fmt"
	"zipline/data"
	"zipline/models"
	"zipline/services"
	"zipline/store"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	// Initialize the inventory service with a memory store
	service := services.NewInventoryService(&store.MemoryStore{})

	// Initialize catalog
	service.InitCatalog(data.ProductList)

	// Simulate a restock
	restock := []models.Item{
		{ProductID: 0, Quantity: 2},
		{ProductID: 10, Quantity: 4},
	}
	service.ProcessRestock(restock)

	// Simulate an order
	order := models.Order{
		OrderID: 123,
		Requested: []models.Item{
			{ProductID: 0, Quantity: 1},
			{ProductID: 10, Quantity: 2},
		},
	}

	if err := service.ProcessOrder(order); err != nil {
		fmt.Println("Failed to process order:", err)
	}
}
