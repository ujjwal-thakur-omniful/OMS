package mongodb

import (
	"context"

	"example.com/m/internal/models"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() (*mongo.Client, error) {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Change to your MongoDB URI
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return nil, err
    }
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        return nil, err
    }
    return client, nil
}

func InsertOrders(ctx context.Context, orders []*models.Order, db *mongo.Client) error {
	collection := db.Database("order").Collection("OMS")


	documents := make([]interface{}, len(orders))
	for i, order := range orders {
		documents[i] = order
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}