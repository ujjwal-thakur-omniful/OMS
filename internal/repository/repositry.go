package repository

import (
	"context"

	"example.com/m/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func NewRepository(db *mongo.Client) *OrderRepository {
	return &OrderRepository{
		db:         db,
		collection: db.Database("order").Collection("OMS"),
	}
}

func InsertOrders(orders []*models.Order, db *mongo.Client) error {
	collection := db.Database("order").Collection("OMS")


	documents := make([]interface{}, len(orders))
	for i, order := range orders {
		documents[i] = order
	}

	_, err := collection.InsertMany(context.Background(), documents)
	return err
}