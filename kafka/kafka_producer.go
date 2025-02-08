package kafka

import (
	"context"
	"log"
	"strconv"

	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
)

type KafkaProducer struct {
	Producer *kafka.ProducerClient
}

var producerInstance *KafkaProducer

func setProducer(producer *kafka.ProducerClient) {
	if producerInstance == nil {
		producerInstance = &KafkaProducer{}
	}
	producerInstance.Producer = producer
}

func getProducer() *kafka.ProducerClient {
	if producerInstance != nil {
		return producerInstance.Producer
	}
	return nil
}

// Initialize Kafka Producer
func InitializeKafkaProducer() {
	producer := kafka.NewProducer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("my-producer"),
		kafka.WithKafkaVersion("2.8.1"),
	)
	log.Printf("Initialized Kafka Producer")
	setProducer(producer)
}

func PublishMessageToKafka(bytesOrderItem []byte, orderID int64) {
	ctx := context.Background()
	msg := &pubsub.Message{
		Topic: "oms-service-topic",
		Key:   strconv.FormatInt(orderID, 10),
		Value: bytesOrderItem,
		Headers: map[string]string{
			"custom-header": "value",
		},
	}

	producer := getProducer()
	err := producer.Publish(ctx, msg)
	if err != nil {
		log.Println("Error publishing message to kafka")
	} else {
		log.Println("Message published to kafka")
	}
}
