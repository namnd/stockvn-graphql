package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	Ctx       context.Context
	Client    *mongo.Client
	Sectors   *mongo.Collection
	Companies *mongo.Collection
	Trades    *mongo.Collection
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
		Ctx:       ctx,
		Client:    client,
		Sectors:   client.Database("stockvn").Collection("sectors"),
		Companies: client.Database("stockvn").Collection("companies"),
		Trades:    client.Database("stockvn").Collection("trades"),
	}
}

func (db mongoDB) Disconnect() {
	db.Client.Disconnect(db.Ctx)
}
