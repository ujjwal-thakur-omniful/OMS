package kafka

import (
	"context"
	"encoding/json"

	"time"

	intersvc "example.com/m/internal/inter_service"
	"example.com/m/internal/models"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

// Implement message handler
type MessageHandler struct{}

// Process implements pubsub.IPubSubMessageHandler.
func (h *MessageHandler) Process(ctx context.Context, message *pubsub.Message) error {
	log.Printf("Received message: %s", string(message.Value))

	// Define a variable to hold the parsed data
	var order models.KafkaResponseOrderMessage
	err := json.Unmarshal(message.Value, &order)
	if err != nil {
		log.WithError(err).Error("Failed to parse Kafka message")
		return err
	}

	// Call WMS Inventory Checking logic from here
	err = intersvc.ValidateInventory(ctx, order)
	if err != nil {
		log.WithError(err).Error("Inventory validation failed \n")
		return err
	}

	return nil
}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	// Process message
	return nil
}

// Initialize Kafka Consumer
func InitializeKafkaConsumer(ctx context.Context) {
	consumer := kafka.NewConsumer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("my-consumer"),
		kafka.WithKafkaVersion("2.8.1"),
		kafka.WithRetryInterval(time.Second),
	)

	handler := &MessageHandler{}
	consumer.RegisterHandler("oms-service-topic", handler)
	consumer.Subscribe(ctx)
}
