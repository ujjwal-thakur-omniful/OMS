package router

import (
	"context"

	controller "example.com/m/internal/controllers"
	"example.com/m/internal/repository"

	mongodb "example.com/m/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
	"go.mongodb.org/mongo-driver/mongo"
)

func Initialize(ctx context.Context, s *http.Server) (err error) {

	// Setup WMS Routes
	omsV1 := s.Engine.Group("/api/v1")
	var db *mongo.Client

	if err != nil {
		return err
	}
	db,err=mongodb.ConnectMongoDB()
	if err != nil {
		return err
	}

	repository.NewRepository(db)


	omsV1.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "mst"})
	})

	omsV1.POST("/create-order", controller.CreateBulkOrders)


	return

}
