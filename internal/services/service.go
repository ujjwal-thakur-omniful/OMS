package service

import (
	"context"

	parse_csv "example.com/m/internal/CSV"
	"example.com/m/internal/domain"
	"example.com/m/internal/models"
	"example.com/m/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"

	error2 "github.com/omniful/go_commons/error"
)

var DB *mongo.Client

type Service struct {
	repository domain.OrderRepository
}

func NewService(repository domain.OrderRepository) domain.OredeService {
	return &Service{repository: repository}
}

func (s *Service) CreateOrder(ctx context.Context, order models.Order) (models.Order, error2.CustomError) {
	// Call the repository to create the order
	createdOrder, cusErr := s.repository.CreateOrder(ctx, order)
	if cusErr.Exists() {

		return models.Order{}, cusErr
	}

	return createdOrder, error2.CustomError{}
}
func CreateBulkOrder(filePath string) error2.CustomError {
	// Logic to process the bulk order from the provided file path
	orders, err := parse_csv.ParseCSV(filePath)
	if err != nil {
		return error2.CustomError{}
	}

	repository.InsertOrders(orders, DB)
	if err != nil {
		return error2.CustomError{}
	}

	return error2.CustomError{}
}
