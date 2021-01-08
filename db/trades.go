package db

import (
	"github.com/namnd/stockvn-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertTrades(trades []*model.Trade) (*mongo.InsertManyResult, error) {
	db := Connect()
	defer db.Disconnect()
	var docs []interface{}
	for _, trade := range trades {
		docs = append(docs, trade)
	}
	return db.Trades.InsertMany(db.Ctx, docs)
}

type tradeNotFoundError struct {
	msg string
}

func (e *tradeNotFoundError) Error() string {
	return e.msg
}

func FindTrades(code string) (trades []*model.Trade, err error) {
	filter := bson.M{"code": code}
	db := Connect()
	options := options.Find()
	options.SetSort(map[string]int{"timestamp": 1})

	cursor, err := db.Trades.Find(db.Ctx, filter, options)
	defer cursor.Close(db.Ctx)
	defer db.Disconnect()
	if err != nil {
		return nil, err
	}
	if err = cursor.All(db.Ctx, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
