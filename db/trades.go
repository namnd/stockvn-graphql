package db

import (
	"fmt"

	"github.com/namnd/stockvn-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func UpdateTrades(trades []*model.Trade) error {
	db := Connect()
	defer db.Disconnect()
	options := options.FindOneAndUpdate().SetUpsert(true)
	for _, trade := range trades {
		filter := bson.M{"code": trade.Code, "date": trade.Date}
		set := bson.M{}
		if trade.OpenPrice > 0 {
			set["openPrice"] = trade.OpenPrice
		}
		if trade.HighPrice > 0 {
			set["highPrice"] = trade.HighPrice
		}
		if trade.LowPrice > 0 {
			set["lowPrice"] = trade.LowPrice
		}
		if trade.AvgPrice > 0 {
			set["avgPrice"] = trade.AvgPrice
		}
		if trade.BuyOrder > 0 {
			set["buyOrder"] = trade.BuyOrder
		}
		if trade.BuyVolume > 0 {
			set["buyVolume"] = trade.BuyVolume
		}
		if trade.SellOrder > 0 {
			set["sellOrder"] = trade.SellOrder
		}
		if trade.SellVolume > 0 {
			set["sellVolume"] = trade.SellVolume
		}
		if trade.MatchedVolume > 0 {
			set["matchedVolume"] = trade.MatchedVolume
		}
		if trade.MatchedValue > 0 {
			set["matchedValue"] = trade.MatchedValue
		}
		if trade.PutThroughOrder > 0 {
			set["putThroughOrder"] = trade.PutThroughOrder
		}
		if trade.PutThroughVolume > 0 {
			set["putThroughVolume"] = trade.PutThroughVolume
		}
		if trade.PutThroughValue > 0 {
			set["putThroughValue"] = trade.PutThroughValue
		}
		if trade.ForeignRemainVolume > 0 {
			set["foreignRemainVolume"] = trade.ForeignRemainVolume
		}
		if trade.ForeignBuyVolume > 0 {
			set["foreignBuyVolume"] = trade.ForeignBuyVolume
		}
		if trade.ForeignBuyValue > 0 {
			set["foreignBuyValue"] = trade.ForeignBuyValue
		}
		if trade.ForeignSellVolume > 0 {
			set["foreignSellVolume"] = trade.ForeignSellVolume
		}
		if trade.ForeignSellValue > 0 {
			set["foreignSellValue"] = trade.ForeignSellValue
		}
		if trade.ForeignPutThroughVolume > 0 {
			set["foreignPutThroughVolume"] = trade.ForeignPutThroughVolume
		}
		if trade.ForeignPutThroughValue > 0 {
			set["foreignPutThroughValue"] = trade.ForeignPutThroughValue
		}
		if len(set) > 0 {
			update := bson.M{
				"$set": set,
			}
			err := db.Trades.FindOneAndUpdate(db.Ctx, filter, update, options).Decode(trade)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

type tradeNotFoundError struct {
	msg string
}

func (e *tradeNotFoundError) Error() string {
	return e.msg
}

type FindTradesOption map[string]string

func FindTrades(filter primitive.M, options *options.FindOptions) (trades []*model.Trade, err error) {
	db := Connect()

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
