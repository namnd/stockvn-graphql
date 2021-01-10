package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/namnd/stockvn-graphql/db"
	"github.com/namnd/stockvn-graphql/graph/generated"
	"github.com/namnd/stockvn-graphql/graph/model"
	"github.com/namnd/stockvn-graphql/scraper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *queryResolver) Sectors(ctx context.Context, exchange *string) ([]*model.Sector, error) {
	return db.FindSectors(exchange)
}

func (r *queryResolver) Companies(ctx context.Context, searchParams *model.CompanySearchParams) ([]*model.Company, error) {
	return db.FindCompanies(searchParams)
}

func (r *queryResolver) Company(ctx context.Context, exchange string, code string) (*model.Company, error) {
	return db.FindCompany(exchange, code)
}

func (r *queryResolver) Trades(ctx context.Context, code string) ([]*model.Trade, error) {
	// By default, get data of the last 5 years
	fromDate := time.Now().AddDate(-5, 0, 0) // 5 years ago
	toDate := time.Now().AddDate(0, 0, -1)   // yesterday

	// Find the latest trade of this stock in the database
	filter := bson.M{"code": code}
	options := options.Find()
	options.SetSort(map[string]int{"date": 1})
	trades, err := db.FindTrades(filter, options)
	if err != nil {
		return nil, err
	}
	if len(trades) > 0 { // If there is any, set fromDate to the latest one
		if lastTrade := trades[len(trades)-1]; lastTrade != nil {
			fromDate = lastTrade.Date.AddDate(0, 0, 1)
		}
	}

	// If yesterday is the latest, skip the scraper
	x := toDate.Sub(fromDate)
	if x.Hours() <= 0 {
		return trades, nil
	}

	// Otherwise, let crawl some new data
	newTrades, err := scraper.GetTrades(code, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	if len(newTrades) > 0 { // If there is any, insert to database
		_, err = db.InsertTrades(newTrades)
		if err != nil {
			return nil, err
		}
	}

	trades = append(trades, newTrades...)
	return trades, nil
}

func (r *queryResolver) TradeMatching(ctx context.Context, code string) ([]*model.Trade, error) {
	filter := bson.M{"code": code, "openPrice": bson.M{"$gt": 0}}
	options := options.Find()
	options.SetSort(map[string]int{"date": 1})
	trades, err := db.FindTrades(filter, options)
	if err != nil {
		return nil, err
	}

	if len(trades) > 0 {
		if lastTrade := trades[len(trades)-1]; lastTrade != nil {
			return trades, nil
		}
	}

	// Let's update matching data to those trades
	newTrades, err := scraper.GetTradeMatchingResults(code)
	if err != nil {
		return nil, err
	}

	err = db.UpdateTrades(newTrades)
	if err != nil {
		return nil, err
	}

	// Update database
	return newTrades, nil
}

func (r *queryResolver) TradePutThrough(ctx context.Context, code string) ([]*model.Trade, error) {
	filter := bson.M{"code": code, "putThroughVolume": bson.M{"$gt": 0}}
	options := options.Find()
	options.SetSort(map[string]int{"date": 1})
	trades, err := db.FindTrades(filter, options)
	if err != nil {
		return nil, err
	}

	if len(trades) > 0 {
		if lastTrade := trades[len(trades)-1]; lastTrade != nil {
			return trades, nil
		}
	}

	newTrades, err := scraper.GetPutThroughResults(code)

	if err != nil {
		return nil, err
	}

	err = db.UpdateTrades(newTrades)
	if err != nil {
		return nil, err
	}

	// Update database
	return newTrades, nil
}

func (r *queryResolver) TradeForeign(ctx context.Context, code string) ([]*model.Trade, error) {
	filter := bson.M{"code": code, "foreignRemainVolume": bson.M{"$gt": 0}}
	options := options.Find()
	options.SetSort(map[string]int{"date": 1})
	trades, err := db.FindTrades(filter, options)
	if err != nil {
		return nil, err
	}

	if len(trades) > 0 {
		if lastTrade := trades[len(trades)-1]; lastTrade != nil {
			return trades, nil
		}
	}
	newTrades, err := scraper.GetForeignResults(code)
	if err != nil {
		return nil, err
	}

	err = db.UpdateTrades(newTrades)
	if err != nil {
		return nil, err
	}

	// Update database
	return newTrades, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
