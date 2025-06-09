package services

import (
	"fmt"
	"zipline/constants"
	"zipline/data"
	"zipline/models"
	"zipline/store"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type (
	// ProductInfo represents a list of products in the catalog.
	ProductInfo []models.Product
	// Restock represents a list of items to be restocked.
	Restock []models.Item
	// Shipment represents the details of a shipment for an order.
	Shipment struct {
		OrderID int           `json:"order_id"`
		Shipped []models.Item `json:"shipped"`
	}
)

// InventoryService provides methods to manage product inventory, process orders, and handle restocks.
type InventoryService struct {
	store *store.MemoryStore
}

// NewInventoryService creates a new instance of InventoryService with the provided memory store.
func NewInventoryService(s *store.MemoryStore) *InventoryService {
	return &InventoryService{store: s}
}

// InitCatalog initializes the product catalog with the provided product information.
func (s *InventoryService) InitCatalog(productInfo ProductInfo) error {
	if len(productInfo) == 0 {
		return fmt.Errorf("catalog initialization failed: product list is empty")
	}

	s.store.InitCatalog(productInfo)
	return nil
}

// GetProductList returns a list of all products in the catalog.
func (s *InventoryService) FindProduct(productID int) (*models.Product, error) {
	product, found := s.store.GetProduct(productID)
	if !found {
		return nil, fmt.Errorf("product with ID %d not found", productID)
	}

	return &product, nil
}

// ProcessOrder processes an order by checking stock availability, preparing shipments, and handling partial shipments.
func (s *InventoryService) ProcessOrder(order models.Order) error {
	if len(order.Requested) == 0 {
		return fmt.Errorf("order %d is empty", order.OrderID)
	}

	var (
		pending         []models.Item
		currentShipment []models.Item
		currentWeight   int
	)

	flushShipment := func() {
		if len(currentShipment) == 0 {
			return
		}
		err := s.ShipPackage(models.Order{
			OrderID:   order.OrderID,
			Requested: currentShipment,
		})
		if err != nil {
			// If shipping fails, requeue everything in currentShipment
			pending = append(pending, currentShipment...)
		}
		currentShipment = nil
		currentWeight = 0
	}

	for _, item := range order.Requested {
		product, err := s.FindProduct(item.ProductID)
		if err != nil {
			return errors.Wrapf(err, "failed to find product with ID %d", item.ProductID)
		}

		available := s.store.GetStock(item.ProductID)
		toShip := min(item.Quantity, available)
		remaining := item.Quantity

		for toShip > 0 {
			if currentWeight+product.MassG > constants.MAX_SHIPMENT_WEIGHT_IN_GRAMS {
				flushShipment()
			}
			if currentWeight+product.MassG <= constants.MAX_SHIPMENT_WEIGHT_IN_GRAMS {
				currentShipment = append(currentShipment, models.Item{
					ProductID: item.ProductID,
					Quantity:  1,
				})
				currentWeight += product.MassG
				toShip--
				remaining--
			}
		}

		// Queue what's not ship-ready
		if remaining > 0 {
			pending = append(pending, models.Item{
				ProductID: item.ProductID,
				Quantity:  remaining,
			})
		}
	}

	flushShipment()

	if len(pending) > 0 {
		s.store.AddPendingOrder(models.Order{
			OrderID:   order.OrderID,
			Requested: pending,
		})
	} else {
		s.store.ClearPendingOrders()
	}

	return nil
}

// ProcessRestock processes a restock request by adding stock for each item in the restock list.
func (s *InventoryService) ProcessRestock(restock Restock) error {
	if len(restock) == 0 {
		return fmt.Errorf("no restock items provided")
	}

	for _, item := range restock {
		if _, found := s.store.GetProduct(item.ProductID); !found {
			log.Warnf("Restock skipped unknown product ID %d", item.ProductID)
			continue
		}
		s.store.AddStock(item.ProductID, item.Quantity)
		log.Infof("Restocked product %d with quantity %d", item.ProductID, item.Quantity)
	}

	// Attempt to reprocess all pending orders
	pending := s.store.GetPendingOrders()
	s.store.ClearPendingOrders()

	for _, order := range pending {
		_ = s.ProcessOrder(order)
	}

	return nil
}

// ShipPackage ships a package for the given order, validating stock and weight limits.
func (s *InventoryService) ShipPackage(order models.Order) error {
	totalWeight := 0

	for _, item := range order.Requested {
		product, found := s.store.GetProduct(item.ProductID)
		if !found {
			log.Errorf("unknown product ID %d during shipment", item.ProductID)
			return fmt.Errorf("unknown product ID %d", item.ProductID)
		}
		stock := s.store.GetStock(item.ProductID)
		if stock < item.Quantity {
			log.Errorf("insufficient stock for product %d (wanted %d, have %d)", item.ProductID, item.Quantity, stock)
			return fmt.Errorf("not enough stock for product %d (wanted %d, have %d)", item.ProductID, item.Quantity, stock)
		}

		totalWeight += product.MassG * item.Quantity
		if totalWeight > constants.MAX_SHIPMENT_WEIGHT_IN_GRAMS {
			log.Errorf("shipment for order %d exceeds weight: %dg > %dg", order.OrderID, totalWeight, constants.MAX_SHIPMENT_WEIGHT_IN_GRAMS)
			return fmt.Errorf("shipment weight exceeds limit (%dg > %dg)", totalWeight, constants.MAX_SHIPMENT_WEIGHT_IN_GRAMS)
		}
	}

	// Deduct stock now that it's valid
	for _, item := range order.Requested {
		s.store.DeductStock(item.ProductID, item.Quantity)
	}

	// Output the shipment details
	log.Infof("Shipped Order #%d", order.OrderID)
	for _, item := range order.Requested {
		log.Infof("  - Product %d x %d", item.ProductID, item.Quantity)
	}
	log.Infof("  Total weight: %dg", totalWeight)

	return nil
}

// PrintProductList prints the list of products in the catalog.
func (s *InventoryService) PrintProductList() {
	for _, product := range data.ProductList {
		fmt.Printf("Product ID: %d, Name: %s, Mass: %dg\n", product.ProductID, product.ProductName, product.MassG)
	}
}
