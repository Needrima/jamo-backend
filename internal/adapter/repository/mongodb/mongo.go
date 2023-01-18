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

func ConnectToMongo(dbUsername string, dbPassword string, dbname string, dbConnString string) (ports.Repository, error) {
	helper.LogEvent("INFO", "Establishing mongoDB connection with given credentials...")

	var client *mongo.Client
	var err error

	if dbConnString == "" {
		credentials := options.Credential{
			Username: dbUsername,
			Password: dbPassword,
		}

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credentials) // Connect to mongodb with credentials
		helper.LogEvent("INFO", "Connecting to MongoDB...")
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			helper.LogEvent("ERROR", "connecting to mongodb")
			return nil, err
		}

		// Check the connection
		helper.LogEvent("INFO", "Confirming MongoDB Connection...")
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			helper.LogEvent("ERROR", "pinging mongodb")
			return nil, err
		}

	} else {
		clientOptions := options.Client().ApplyURI(dbConnString)
		helper.LogEvent("INFO", "Connecting to MongoDB...")
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			helper.LogEvent("ERROR", "connecting to mongodb")
			return nil, err
		}

		// Check the connection
		helper.LogEvent("INFO", "Confirming MongoDB Connection...")
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			helper.LogEvent("ERROR", "pinging mongodb")
			return nil, err
		}
	}

	//helper.LogEvent("Info", "Connected to MongoDB!")
	helper.LogEvent("INFO", "Establishing Database collections and indexes...")
	conn := client.Database(dbname)

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
