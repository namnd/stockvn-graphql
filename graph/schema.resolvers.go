package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/namnd/stockvn-graphql/db"
	"github.com/namnd/stockvn-graphql/graph/generated"
	"github.com/namnd/stockvn-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sectors(ctx context.Context) ([]*model.Sector, error) {
	var sectors []*model.Sector
	cursor, err := db.Connect().Sectors.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &sectors); err != nil {
		log.Fatal(err)
	}
	return sectors, nil
}

func (r *queryResolver) Companies(ctx context.Context) ([]*model.Company, error) {
	var companies []*model.Company
	cursor, err := db.Connect().Companies.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &companies); err != nil {
		log.Fatal(err)
	}
	return companies, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
