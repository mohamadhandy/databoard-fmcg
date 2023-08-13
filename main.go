package main

import (
	"flag"
	"klikdaily-databoard/service"
	"log"
)

func main() {
	isIndexerService := flag.Bool("indexer", false, "Run as indexer service")
	isHttpService := flag.Bool("http", true, "Run as http service")
	flag.Parse()

	// router := gin.Default()

	// router.POST("/users", authMiddleware, userController.RegisterUser)
	// router.POST("/login", userController.LoginUser)
	// router.GET("/users", authMiddleware, userController.GetAllUsers)
	// router.GET("/users/:user_id", authMiddleware, userController.GetUsersByID)
	// router.PUT("/users/:user_id", authMiddleware, userController.UpdateUser)
	// router.DELETE("/users/:user_id", authMiddleware, userController.DeleteUser)
	// es, err := config.SetupElasticsearch()
	// if err != nil {
	// 	panic(err)
	// }
	// // Create the Elasticsearch index
	// indexName := "elasticsearch_index"
	// // Check if the index exists
	// exists, err := config.IndexExists(es, indexName)
	// if err != nil {
	// 	log.Fatalf("Error checking index existence: %s", err)
	// }

	// if !exists {
	// 	// Create the Elasticsearch index
	// 	err = config.CreateIndex(es, indexName)
	// 	if err != nil {
	// 		log.Fatalf("Error creating index: %s", err)
	// 	}
	// 	log.Printf("Index %s created", indexName)
	// } else {
	// 	log.Printf("Index %s already exists", indexName)
	// }
	// rabbitMQ, errRabbitMQ := config.SetupRabbitMQ()
	// if errRabbitMQ != nil {
	// 	panic(errRabbitMQ)
	// }
	// routes.RouteAPI(router, context.Background(), config.NewConnection(), config.NewConnectionRedis(), es)
	// router.Run("localhost:9000")

	if *isIndexerService {
		errC, err := service.LoadIndexerService()
		if err != nil {
			log.Fatalf("Couldn't run: %s", err)
		}

		if err := <-errC; err != nil {
			log.Fatalf("Error while running: %s", err)
		}
	}

	if !*isIndexerService && *isHttpService {
		service.LoadHttpService()
		return
	}
}
