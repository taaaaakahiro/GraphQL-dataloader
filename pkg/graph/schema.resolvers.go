package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
)

func (r *mutationResolver) CreateMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Messages(ctx context.Context, userID string) ([]*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllMessages(ctx context.Context) ([]*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
