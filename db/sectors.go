package db

import (
	"fmt"

	"github.com/namnd/stockvn-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func FindSectors(exchange *string) ([]*model.Sector, error) {
	filter := bson.M{}
	if exchange != nil && *exchange != "all" {
		filter = bson.M{"exchange": *exchange}
	}
	db := Connect()
	fmt.Println(exchange)
	cursor, err := db.Sectors.Find(db.Ctx, filter)
	defer cursor.Close(db.Ctx)
	defer db.Disconnect()

	if err != nil {
		return nil, err
	}
	var sectors []*model.Sector
	if err = cursor.All(db.Ctx, &sectors); err != nil {
		return nil, err
	}
	fmt.Println(sectors)
	return sectors, nil
}
