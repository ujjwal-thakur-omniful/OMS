package appinit

import (
	"context"
	"fmt"

	
	"github.com/omniful/go_commons/sqs"
)


var exportedPublisher *sqs.Publisher
var exportedQueue *sqs.Queue

func GetNewSQSPublisher() *sqs.Publisher {
	return exportedPublisher
}

func setNewSQSPublisher(publisher *sqs.Publisher) {
	exportedPublisher = publisher
}

func getSQSUrl(ctx context.Context) *string {
	sqsURL, err := sqs.GetUrl(ctx, sqs.GetSQSConfig(ctx, false, "", "eu-north-1", "794038235002", "https://sqs.eu-north-1.amazonaws.com/"), "bulk_order")
	fmt.Println("sqsURL: ", sqsURL)
	if err != nil {
		fmt.Println("error in connecting sqs")
	}
	return sqsURL
}

func initialiseSQSProducer(ctx context.Context) {

	sqsURL := getSQSUrl(ctx)
	fmt.Println("Queue URL: ", *sqsURL)

	newQueue, err := sqs.NewStandardQueue(ctx, "bulk_order", sqs.GetSQSConfig(ctx, false, "", "eu-north-1", "794038235002", "https://sqs.eu-north-1.amazonaws.com/"))
	if err != nil {
		fmt.Println("Error in creating queue")
	}
	exportedQueue = newQueue
	NewSQSPublisher := sqs.NewPublisher(newQueue)
	setNewSQSPublisher(NewSQSPublisher)
}
