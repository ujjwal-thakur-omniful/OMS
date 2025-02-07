package appinit

import (
	"context"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
)
var DB *mongo.Client
func Initialize(ctx context.Context) {
	initializeLog(ctx)
	initialiseSQSProducer(ctx)
	StartConsumerWorker(ctx)
	

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
