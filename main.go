package main

import (
	"context"

	"time"

	appinit "example.com/m/init"
	"example.com/m/router"
	"github.com/omniful/go_commons/config"

	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/shutdown"
	
)

const (
// modeWorker     = "worker"
// modeHttp       = "http"
// modeMigration  = "migration"
// upMigration    = "up"
// downMigration  = "down"
// forceMigration = "force"
)

func main() {

	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initialising config, err: %v", err)
		panic(err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Error while getting context from config, err: %v", err)
		panic(err)
	}
	
	//initoialise connection
	appinit.Initialize(ctx)

	//run server
	//runHttpServer(ctx)

	runHttpServer(ctx)

}
func runHttpServer(ctx context.Context) {

	server := http.InitializeServer(config.GetString(ctx, "server.port"), 10*time.Second, 10*time.Second, 70*time.Second)

	// Initialize middlewares and routes
	err := router.Initialize(ctx, server)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	//

	log.Infof("Starting server on port" + config.GetString(ctx, "server.port"))

	err = server.StartServer("OMS-service")
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	<-shutdown.GetWaitChannel()
}
// func ConnectMongoDB() (*mongo.Client, error) {
//     clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Change to your MongoDB URI
//     client, err := mongo.Connect(context.TODO(), clientOptions)
//     if err != nil {
//         return nil, err
//     }
//     err = client.Ping(context.TODO(), nil)
//     if err != nil {
//         return nil, err
//     }
//     return client, nil
// }