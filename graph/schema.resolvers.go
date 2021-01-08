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
	// Find the latest trade of this stock in the database
	trades, err := db.FindTrades(code)
	if err != nil {
		return nil, err
	}
	fromDate := time.Now().AddDate(-5, 0, 0) // 5 years ago
	toDate := time.Now().AddDate(0, 0, -1)   // yesterday
	if len(trades) > 0 {
		if lastTrade := trades[len(trades)-1]; lastTrade != nil {
			fromDate = time.Unix(int64(lastTrade.Timestamp)/1000, 0)
		}
	}
	x := toDate.Sub(fromDate)
	if x.Hours() <= 0 {
		return trades, nil
	}

	newTrades, err := scraper.GetTrades(code, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	if len(newTrades) > 0 {
		_, err = db.InsertTrades(newTrades)
		if err != nil {
			return nil, err
		}
	}

	trades = append(trades, newTrades...)
	return trades, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
