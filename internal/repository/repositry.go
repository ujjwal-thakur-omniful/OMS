package repository

import (
	"context"
	"fmt"

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
		collection: db.Database("order-oms").Collection("OMS"),
	}
}

func InsertOrders(ctx context.Context, orders []*models.Order, db *mongo.Client) error {
	if len(orders) == 0 {
		fmt.Println("No orders to insert")
		return nil // No orders to insert, so return without an error
	}

	collection := db.Database("order-oms").Collection("OMS")

	// Convert orders from []*models.Order to []interface{}
	documents := make([]interface{}, len(orders))
	for i, order := range orders {
		documents[i] = *order // Dereferencing pointer
	}

	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		fmt.Println("Error inserting orders:", err)
		return err
	}

	fmt.Println("Orders inserted successfully")
	return nil
}
