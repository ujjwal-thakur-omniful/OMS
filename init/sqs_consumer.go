// package appinit

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	// "oms-service/intersvc"

// 	parse_csv "example.com/m/internal/CSV"
// 	"example.com/m/internal/models"
// 	"example.com/m/internal/repository"

// 	// "example.com/m/internal/parse_csv"
// 	// "example.com/m/internal/repository"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/sqs"
// 	"github.com/omniful/go_commons/log"
// )

// func getConfig() *aws.Config {
// 	if awsConfig == nil {
// 		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-north-1"))
// 		if err != nil {
// 			panic("Unable to connect to AWS")
// 		}
// 		awsConfig = &cfg
// 		return awsConfig
// 	}
// 	return awsConfig
// }

// func initialiseSQSConsumer(ctx context.Context) {

// 	sqsClient := sqs.NewFromConfig(*getConfig())

// 	sqURL := getSQSUrl(ctx)
// 	fmt.Println("Queue URL: ", *sqURL)

// 	// This will constantly listen to the SQS queue and print the messages
// 	// for {
// 	messagesResult, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
// 		QueueUrl: sqURL,
// 	})
// 	if err != nil {
// 		fmt.Println("Unable to receive mesaages from SQS: ", err)
// 	}

// 	var sqsMessage models.SQSMessage
// 	for _, message := range messagesResult.Messages {
// 		fmt.Println("Message: ", *message.Body)
// 		if err := json.Unmarshal([]byte(*message.Body), &sqsMessage); err != nil {
// 			log.Printf("Error unmarshaling message: %v", err)
// 			continue
// 		}
// 	}
// 	orders, err := parse_csv.ParseCSV(sqsMessage.FilePath)
// 	if err != nil {
// 		log.Printf("Error parsing CSV: %v", err)
// 	}
// 	fmt.Println("Orders: ", orders)
// 	// Parse CSV
// 	// orders, err := parse_csv.ParseCSV(sqsMessage.FilePath)
// 	// if err != nil {
// 	// 	log.Printf("Error parsing CSV: %v", err)
// 	// }
// 	// fmt.Println("Orders: ", orders)
// 	// Save to MongoDB
// 	if err := repository.InsertOrders(ctx, orders, DB); err != nil {
// 		log.Printf("Error saving orders to database: %v", err)
// 	}

//		// Interservice Call to WMS to check
//	}
package appinit

import (
	// "OMS/service"
	"context"
	"encoding/json"
	"fmt"
	"log"

	service "example.com/m/internal/services"
	"github.com/omniful/go_commons/sqs"
)

type ExampleHandler struct{}

func (h *ExampleHandler) Process(ctx context.Context, message *[]sqs.Message) error {
	//TODO implement me
	//panic("implement me")

	for _, msg := range *message {
		err := h.Handle(&msg)
		if err != nil {

		}
	}
	return nil
}

func (h *ExampleHandler) Handle(msg *sqs.Message) error {
	fmt.Println("Processing message:", string(msg.Value))
	var event struct {
		FilePath string `json:"filePath"`
	}

	if err := json.Unmarshal([]byte(msg.Value), &event); err != nil {
		log.Printf("Failed to parse SQS message: %v", err)
		return err
	}

	// Call service to process the bulk order
	err := service.CreateBulkOrder(event.FilePath)
	if err.Exists() {
		log.Printf("Failed to process bulk order: %v", err)
		return err
	}

	if err := exportedQueue.Remove(msg.ReceiptHandle); err != nil {
		log.Printf("Failed to remove message from queue: %v", err)
		return err
	}


	fmt.Println("Bulk order processing complete")

	return nil
}

func StartConsumerWorker(ctx context.Context) {

	// Set up consumer
	handler := &ExampleHandler{}
	consumer, err := sqs.NewConsumer(
		exportedQueue,
		1, // Number of workers
		1, // Concurrency per worker
		handler,
		10, // Max messages count
		30, // Visibility timeout

		false, // Is async
		false, // Send batch message
	)

	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	consumer.Start(ctx)


}
