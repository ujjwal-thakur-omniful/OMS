
package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/m/internal/models"

	"github.com/omniful/go_commons/http"

	interservice_client "github.com/omniful/go_commons/interservice-client"
)

func ValidateInventory(ctx context.Context, order models.KafkaResponseOrderMessage) error {

	log.Printf("Validating inventory for order ID: %s \n", order.OrderID)

	config := interservice_client.Config{
		ServiceName: "oms-service",
		BaseURL:     "http://localhost:8080/api/v1/order",
		Timeout:     5 * time.Second,
	}

	client, err := interservice_client.NewClientWithConfig(config)
	if err != nil {
		return err
	}

	url := config.BaseURL + "validate_inventory"
	bodyBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req := &http.Request{
		Url:     url, // Use configured URL
		Body:    bytes.NewReader(bodyBytes),
		Timeout: 7 * time.Second,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	_, intersvcErr := client.Post(req, "/")
	if intersvcErr.StatusCode.Is4xx() {
		fmt.Println("inventory validation failed after  interservice call to wms-service ")
	} else {
		fmt.Println("Inventory validation")
	}

	return nil
}
