package db

import (
	"context"
	"log"
	"time"

	"github.com/namnd/stockvn-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Trade struct {
	Code       string `bson:"code"`
	Timestamp  int    `bson:"timestamp"`
	ClosePrice int    `bson:"close_price"`
	Volume     int    `bson:"volume"`
}

func InsertTrade(code string, trade *model.Trade) (*mongo.InsertOneResult, error) {
	doc := Trade{
		Code:       code,
		Timestamp:  trade.Timestamp,
		ClosePrice: trade.ClosePrice,
		Volume:     trade.Volume,
	}
	db := Connect()
	defer db.Disconnect()
	return db.Trades.InsertOne(db.Ctx, doc)
}

func FindTrades(searchParams *model.TradeSearchParams) map[int]bool {
	layout := "2006-02-01"
	from, err := time.Parse(layout, searchParams.From)
	to, err := time.Parse(layout, searchParams.To)
	filter := bson.M{"code": searchParams.Code}
	timestampQuery := []bson.M{}
	timestampQuery = append(timestampQuery, bson.M{"timestamp": bson.M{"$lte": to.Unix() * 1000}})
	filter["$and"] = append(timestampQuery, bson.M{"timestamp": bson.M{"$gte": from.Unix() * 1000}})
	db := Connect()
	cursor, err := db.Trades.Find(db.Ctx, filter)
	defer cursor.Close(db.Ctx)
	defer db.Disconnect()

	if err != nil {
		log.Println(err)
	}
	var trades []Trade
	if err = cursor.All(context.Background(), &trades); err != nil {
		log.Println(err)
	}
	foundTrades := make(map[int]bool, len(trades))
	for _, trade := range trades {
		foundTrades[trade.Timestamp] = true
	}
	return foundTrades
}
