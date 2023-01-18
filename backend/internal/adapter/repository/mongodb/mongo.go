package repository

import (
	"context"
	"errors"
	"jamo/backend/internal/core/helper"
	ports "jamo/backend/internal/port"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func ConnectToMongo(dbUsername string, dbPassword string, dbname string, dbPort string) (ports.Repository, error) {
	helper.LogEvent("INFO", "Establishing mongoDB connection with given credentials...")

	credentials := options.Credential{
		Username: dbUsername,
		Password: dbPassword,
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:" + dbPort).SetAuth(credentials) // Connect to mongodb with credentials
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:"+dbPort)// Connect to mongodb without credentials
	helper.LogEvent("INFO", "Connecting to MongoDB...")
	db, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		helper.LogEvent("ERROR", "connecting to mongodb")
		return nil, err
	}

	// Check the connection
	helper.LogEvent("INFO", "Confirming MongoDB Connection...")
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		//log.Println(err)
		//log.Fatal(err)
		helper.LogEvent("ERROR", "pinging mongodb")
		return nil, err
	}

	//helper.LogEvent("Info", "Connected to MongoDB!")
	helper.LogEvent("INFO", "Establishing Database collections and indexes...")
	conn := db.Database(dbname)

	productCollection := conn.Collection("products")
	newsletterCollection := conn.Collection("newsletter")
	orderCollection := conn.Collection("orders")
	messagesCollection := conn.Collection("messages")

	return NewInfra(productCollection, newsletterCollection, orderCollection, messagesCollection), nil
}

func GetPage(page string) (*options.FindOptions, error) {
	var limit, e = strconv.ParseInt(helper.Config.PageLimit, 10, 64)
	var pageSize, ee = strconv.ParseInt(page, 10, 64)
	if e != nil || ee != nil {
		return nil, errors.New("page limit or page size")
	}
	findOptions := options.Find().SetLimit(limit).SetSkip(limit * (pageSize - 1))
	return findOptions, nil
}
