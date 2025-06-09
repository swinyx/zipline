# Zipline Inventory & Order Management System

This is a lightweight backend system simulating the core logistics functions of [Zipline](https://flyzipline.com), designed to:

* Track product inventory
* Process incoming restocks
* Handle hospital orders
* Automatically trigger package shipments (not exceeding 1.8kg per shipment)

---

## Features

* Initialize a product catalog
* Handle restocks and update inventory
* Process and split orders into shipment-sized batches
* Defer orders if insufficient stock
* Print `ship_package` instructions to the console

---

## Project Structure

```
zipline/
├── main.go                         # Entry point
├── constants/                      # Project-wide constants
│   └── constants.go
├── models/                         # Shared data types: Product, Item, Order
│   └── models.go
├── services/                       # Core business logic
│   ├── inventory_service.go
│   └── inventory_service_test.go   # Unit tests for InventoryService
├── store/                          # Interfaces and memory-backed store implementations
│   ├── interface.go                # Interfaces: ProductStore, InventoryStore, OrderStore
│   ├── store.go                    # In-memory store implementation
│   └── store_test.go               # Unit tests for the store layer
├── go.mod
└── README.md                       # You're here
```

---

## Getting Started

### Prerequisites

* Go 1.20+
* Git (for cloning)

### Setup

```bash
# Clone the repo
git clone https://github.com/swinyx/zipline.git
cd zipline

# Initialize Go modules (if not already)
go mod tidy

# Run the app
go run main.go
```

---

## Example

In `main.go`:

```go
catalog := []models.Product{
  {ProductID: 0, ProductName: "RBC A+ Adult", MassG: 700},
  {ProductID: 10, ProductName: "FFP A+", MassG: 300},
}
service.InitCatalog(catalog)

restock := []models.Item{
  {ProductID: 0, Quantity: 2},
  {ProductID: 10, Quantity: 4},
}
service.ProcessRestock(restock)

order := models.Order{
  OrderID: 123,
  Requested: []models.Item{
    {ProductID: 0, Quantity: 1},
    {ProductID: 10, Quantity: 3},
  },
}
service.ProcessOrder(order)
```

Expected console output:

```json
{
  "order_id": 123,
  "requested": [
    {"product_id": 0, "quantity": 1},
    {"product_id": 10, "quantity": 3}
  ]
}
```

---

## Running Tests

The project uses Go's built-in testing framework. To run the tests:

```bash
go test ./...
```

This will recursively run all tests across the project.

To run tests with verbose output:

```bash
go test -v ./...
```

If you have individual test files in a folder (e.g., `services/`), you can run them like this:

```bash
go test ./services
```

---

## Architecture Overview

The system is split into well-defined layers:

* `models/`: Basic data types shared across the app (`Product`, `Item`, `Order`).
* `store/`: Contains interfaces for product, inventory, and order storage — and in-memory implementations for fast prototyping and testability.
* `services/`: Business logic lives here. It coordinates inventory updates, shipment logic, and order fulfillment.
* `constants/`: Centralizes business rules like the maximum shipment weight.

This design enables loose coupling between data, logic, and behavior — which improves testability and flexibility.

---

## Business Rules Covered

* Shipments must not exceed 1.8 kg (`1800g`).
* Orders can be partially fulfilled if stock is insufficient.
* Orders not immediately fulfillable are deferred until stock is available.
* All shipments are logged to the console as if being sent to a UI/fulfillment operator.

---

## Key Test Suggestions

You can write unit tests that verify:

* Inventory increases with restock
* Orders that exceed inventory get deferred
* `ship_package` never exceeds 1800g total
* Partial shipments happen correctly in multiple calls
* Edge case: product not found in catalog → handled gracefully

---

## Concurrency Considerations

This prototype does not yet implement concurrency protections (e.g., mutexes).
However, in a real-world setting with multiple goroutines or network calls, shared resources like inventory should be guarded with:

* `sync.Mutex` or `sync.RWMutex`
* Atomic counters for performance-critical paths
* Channel-based work queues or goroutines for batch processing

---

## Future Enhancements

* REST or gRPC API integration
* Persistent storage (e.g., PostgreSQL)
* Dashboard for fulfillment ops
* Metrics and alerting for stockouts

---

## Author

**Yannick Yeboue**
[GitHub](https://github.com/swinyx) • [LinkedIn](https://www.linkedin.com/in/yannick-yeboue/)

---

## License

MIT – feel free to use and modify!
