package appinit

import (
	"context"
	"fmt"
	"time"

	"example.com/m/kafka"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func Initialize(ctx context.Context) {
	initializeLog(ctx)
	initialiseSQSProducer(ctx)
	initializeDB(ctx)
	startSQSConsumer(ctx)
	kafka.InitializeKafkaProducer()
	go kafka.InitializeKafkaConsumer(ctx)


}

// Initialize logging
func initializeLog(ctx context.Context) {
	err := log.InitializeLogger(
		log.Formatter(config.GetString(ctx, "log.format")),
		log.Level(config.GetString(ctx, "log.level")),
	)
	if err != nil {
		log.WithError(err).Panic("unable to initialise log")
	}

}
func initializeDB(c context.Context) {
	fmt.Println("Connecting to mongo...")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://ujjwal3112:ljYWuz0jE73xIAW5@cluster0.tyixv4j.mongodb.net/order-oms")

	var err error
	DB, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	err = DB.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Successfully connected to MongoDB!")
}

// Initialize Redis
// func initializeMongoDB(ctx context.Context) {
// 	uri := config.GetString(ctx, "mongo.uri")
// 	dbName := config.GetString(ctx, "mongo.database")
// 	collectionName := config.GetString(ctx, "mongo.collection")

// 	_, err := mongo.NewMongoDB(uri, dbName, collectionName)
// 	if err != nil {
// 		log.WithError(err).Panic("unable to initialize MongoDB")
// 	}

// }
