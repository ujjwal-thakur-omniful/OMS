package parse_csv

import (
	"context"

	"fmt"
	"log"

	"os"
	"strconv"
	"time"

	validate "example.com/m/internal/events"
	"example.com/m/internal/models"
	"github.com/omniful/go_commons/csv"
)

func ParseCSV(filePath string) ([]*models.Order, error) {
	fmt.Println("Parse CSV function called successfull!")
	fmt.Println("This is the file beign opened: ", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error in opening file path.")
	}
	defer file.Close()

	// Map to group items by order_no
	orderGroups := make(map[string]models.Order)

	// Initialize the CSV reader (based on your previous implementation)
	CSV, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(filePath),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}

	err = CSV.InitializeReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}

	// Process the records and group them by order_no and customer_name
	for !CSV.IsEOF() {
		var records csv.Records
		records, err := CSV.ReadNextBatch()
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Processing records:")
		fmt.Println(records)
		for _, record := range records {
			orderID, _ := strconv.Atoi(record[0])

			sellerID, _ := strconv.Atoi(record[1])
			tenantID, _ := strconv.Atoi(record[2])
			hubID, _ := strconv.Atoi(record[3])

			orderStatus := record[6]
			modeOfPayment := record[7]

			// hub_id
			skuID, _ := strconv.Atoi(record[0])       // Assuming SKU ID is in the 9th column (index 8)
			IntQuantity, _ := strconv.Atoi(record[1]) // Assuming Quantity is in the 8th column (index 7)

			// Convert quantity to integer

			// Check if the  group forderor this order_id already exists
			orderKey := strconv.Itoa(orderID)
			order, exists := orderGroups[orderKey]
			if !exists {
				// If order doesn't exist, create a new order
				now := time.Now()

				order = models.Order{
					// SellerID:     sellerID,
					// HubID:        hubID,
					ID:            orderID,
					SellerID:      sellerID,
					TenantID:      tenantID,
					HubID:         hubID,
					CreatedAt:     now,
					UpdatedAt:     now,
					OrderStatus:   orderStatus,
					ModeOfPayment: modeOfPayment,
					OrderItems:    []models.OrderItem{}, // Start with an empty slice of items
				}
				// Add the new order to the map
				orderGroups[orderKey] = order
			}

			// Create a new OrderItem and append it to the order's OrderItems
			orderItem := models.OrderItem{

				SKUID:    skuID,
				Quantity: IntQuantity,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}
	}
	fmt.Println("orderGroups", orderGroups)
	var orders []*models.Order
	for _, order := range orderGroups {
		orders = append(orders, &order)
	}
	fmt.Println("orders", orders)

	fmt.Println("Final orders:")
	for _, order := range orders {
		fmt.Printf("Order No: %s, Total Items: %d\n", order.ID, len(order.OrderItems))
		go validate.ValidateOrders(order)
	}

	// var orders []models.Order // Use a slice of values instead of pointers
	// for _, order := range orderGroups {
	// 	orders = append(orders, *order) // Dereference the pointer before appending
	// }

	// // Convert to JSON
	// jsonData, err := json.Marshal(orders)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("jsonData",jsonData)
	// fmt.Println("Final orders:")
	// for _, order := range orders {
	// 	fmt.Printf("Order No: %s, Total Items: %d\n", order.ID, len(order.OrderItems))
	// 	// go intersvc.ValidateOrders(order)
	// }

	return orders, nil

}
