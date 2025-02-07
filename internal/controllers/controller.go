package controller

import (
	"encoding/json"

	"os"

	appinit "example.com/m/init"
	"example.com/m/internal/domain"
	"example.com/m/internal/models"
	"github.com/gin-gonic/gin"
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
