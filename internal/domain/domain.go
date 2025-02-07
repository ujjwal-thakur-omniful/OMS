package domain

import (
	"context"

	"example.com/m/internal/models"
	error2 "github.com/omniful/go_commons/error"
)



type OredeService interface {
	// GetOrderDetails(c context.Context,  order_id uint64) (models.Order, error2.CustomError)
	CreateOrder(ctx context.Context, order models.Order) (models.Order,  error2.CustomError)



}
type OrderRepository interface {
	// GetOrder(c context.Context,  order_id uint64) (models.Order, error2.CustomError)
	CreateOrder(ctx context.Context, order models.Order) (models.Order,  error2.CustomError)
}

