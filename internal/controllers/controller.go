package controller

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"

	"os"
	"strconv"

	"time"

	appinit "example.com/m/init"
	"example.com/m/internal/domain"
	"example.com/m/internal/models"
	oms_kafka "example.com/m/kafka"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"

	interservice_client "github.com/omniful/go_commons/interservice-client"
	"github.com/omniful/go_commons/sqs"
)

type Controller struct {
	service domain.OredeService
}

func NewController(service domain.OredeService) *Controller {
	return &Controller{service: service}
}

func CreateBulkOrders(ctx *gin.Context) {
	var request models.BulkOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate if file exists
	if _, err := os.Stat(request.FilePath); os.IsNotExist(err) {
		ctx.JSON(400, gin.H{"error": "File does not exist"})
		return
	}

	// Convert request to bytes
	messageBytes, err := json.Marshal(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to marshal request"})
		return
	}

	newMessage := &sqs.Message{
		GroupId:       "group-323",
		Value:         messageBytes,
		ReceiptHandle: "receipt-abc",
		Attributes:    map[string]string{"key1": "value1", "key2": "value2"},
	}

	// Publish message to SQS
	publisher := appinit.GetNewSQSPublisher()
	if err := publisher.Publish(ctx, newMessage); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to publish message to queue"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Bulk order request queued successfully"})
}



func ValidateOrders(order *models.Order) {
	fmt.Println("Validate fxn called")

	for _, orderItem := range order.OrderItems{
		config := interservice_client.Config{
			ServiceName: "oms-service",
			BaseURL:     "http://localhost:8080/api/v1/order",
			Timeout:     5 * time.Second,
		}

		client, err := interservice_client.NewClientWithConfig(config)
		if err != nil {
			return
		}

		url := config.BaseURL + "validate_order"
		body := map[string]string{
			"sku_id": strconv.Itoa(orderItem.SKUID),
			"hub_id": strconv.Itoa(order.HubID),
		}


		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return
		}

		req := &http.Request{
			Url:     url,
			Body:    bytes.NewReader(bodyBytes),
			Timeout: 7 * time.Second,
			Headers: map[string][]string{


				"Content-Type": {"application/json"},
			},
		}

		resp, intersvcErr := client.Post(req, "/")
		if intersvcErr.StatusCode.Is4xx() {
			fmt.Printf("Error making GET request to validate SKU: %v\n", err)
			return
		} else {
			fmt.Print(resp)
			log.Printf("Order with Order ID: %v having product %v from hub %v is VALID \n", order.ID, orderItem.SKUID, order.HubID)


			// Publish This Order Item in a message to Kafka
			bytesOrderItem, _ := json.Marshal(orderItem)
			oms_kafka.PublishMessageToKafka(bytesOrderItem, int64(order.ID))

			return
		}
	}
}
