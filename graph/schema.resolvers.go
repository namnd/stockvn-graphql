package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/namnd/stockvn-graphql/db"
	"github.com/namnd/stockvn-graphql/graph/generated"
	"github.com/namnd/stockvn-graphql/graph/model"
	"github.com/namnd/stockvn-graphql/scraper"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sectors(ctx context.Context, exchange *string) ([]*model.Sector, error) {
	return db.FindSectors(exchange)
}

func (r *queryResolver) Companies(ctx context.Context, searchParams *model.CompanySearchParams) ([]*model.Company, error) {
	return db.FindCompanies(searchParams)
}

func (r *queryResolver) Company(ctx context.Context, exchange string, code string) (*model.Company, error) {
	return db.FindCompany(exchange, code)
}

func (r *queryResolver) Trades(ctx context.Context, searchParams *model.TradeSearchParams) ([]*model.Trade, error) {
	// First crawl some data
	var trades []*model.Trade
	trades, err := scraper.GetTrades(searchParams)
	fmt.Printf("Crawl %d trades", len(trades))
	if err != nil {
		return nil, err
	}

	// Then check if we already have some of them in database
	oldTrades := db.FindTrades(searchParams)
	fmt.Printf("Found %d trades in db", len(oldTrades))

	// Then save only new data to database
	for _, trade := range trades {
		if _, ok := oldTrades[trade.Timestamp]; !ok {
			_, err := db.InsertTrade(searchParams.Code, trade)
			if err != nil {
				return nil, err
			}
		}
	}
	return trades, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
