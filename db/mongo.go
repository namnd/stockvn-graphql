package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	Client    *mongo.Client
	Companies *mongo.Collection
	Sectors   *mongo.Collection
}

func Connect() mongoDB {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return mongoDB{
		Client:    client,
		Companies: client.Database("stockvn").Collection("companies"),
		Sectors:   client.Database("stockvn").Collection("sectors"),
	}
}

func (db mongoDB) Disconnect() {
	err := db.Client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
