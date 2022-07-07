package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/robertsmoto/skustor/graph/generated"
	"github.com/robertsmoto/skustor/graph/model"
)

func (r *collectionResolver) ImageIds(ctx context.Context, obj *model.Collection) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *collectionResolver) ItemIds(ctx context.Context, obj *model.Collection) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddCollection(ctx context.Context, input model.NewCollection) (*model.Collection, error) {
	collection := &model.Collection{
		ID:       input.ID,
		Document: input.Document,
	}
	r.collections = append(r.collections, collection)
	return collection, nil
}

func (r *mutationResolver) AddCollections(ctx context.Context, input []*model.NewCollection) ([]*model.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collections(ctx context.Context) ([]*model.Collection, error) {
	return r.collections, nil
}

// Collection returns generated.CollectionResolver implementation.
func (r *Resolver) Collection() generated.CollectionResolver { return &collectionResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type collectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
