package data

import "zipline/models"

var (
	// ProductList contains a predefined list of products available in the system.
	ProductList = []models.Product{
		{ProductID: 0, ProductName: "RBC A+ Adult", MassG: 700},
		{ProductID: 1, ProductName: "RBC B+ Adult", MassG: 700},
		{ProductID: 2, ProductName: "RBC AB+ Adult", MassG: 750},
		{ProductID: 3, ProductName: "RBC O- Adult", MassG: 680},
		{ProductID: 4, ProductName: "RBC A+ Child", MassG: 350},
		{ProductID: 5, ProductName: "RBC AB+ Child", MassG: 200},
		{ProductID: 6, ProductName: "PLT AB+", MassG: 120},
		{ProductID: 7, ProductName: "PLT O+", MassG: 80},
		{ProductID: 8, ProductName: "CRYO A+", MassG: 40},
		{ProductID: 9, ProductName: "CRYO AB+", MassG: 80},
		{ProductID: 10, ProductName: "FFP A+", MassG: 300},
		{ProductID: 11, ProductName: "FFP B+", MassG: 300},
		{ProductID: 12, ProductName: "FFP AB+", MassG: 300},
	}
)
