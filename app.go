package main

import (
	"fmt"
	mongoRepository "jamo/backend/internal/adapter/repository/mongodb"
	"jamo/backend/internal/adapter/routes"
	"jamo/backend/internal/core/helper"
	ports "jamo/backend/internal/port"
	"log"
	"os"
)

func main() {
	//Initialize request Log
	helper.InitializeLog()
	//Start DB Connection
	mongoRepo := startMongo()
	helper.LogEvent("INFO", "MongoDB Initialized!")

	// helper.LogEvent("INFO", "Redis Initialized!")

	//Set up routes
	router := routes.SetupRouter(mongoRepo)
	//Print custom message for server start

	helper.LogEvent("INFO", "server started")
	//start server
	port := os.Getenv("PORT")
	if port == "" {
		port = helper.Config.ServicePort
	}
	_ = router.Run(":" + port)
	//api.SetConfiguration
}

func startMongo() ports.Repository {
	helper.LogEvent("INFO", "Initializing Mongo!")
	mongoRepo, err := mongoRepository.ConnectToMongo(
		helper.Config.MongoDbUserName,
		helper.Config.MongoDbPassword,
		helper.Config.MongoDbName,
		os.Getenv("mongo_db_conn_string"),
	)

	if err != nil {
		fmt.Println(err)
		helper.LogEvent("ERROR", "MongoDB database Connection Error: "+err.Error())
		log.Fatal()
	}
	return mongoRepo
}
