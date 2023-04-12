package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.29

import (
	"context"
	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/graph/model"
	"dreamkast-weaver/internal/sqlhelper"
	"log"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, input model.NewVote) (*bool, error) {
	ctx = sqlhelper.WithSqlHelper(ctx, r.sqlHelper)

	if err := r.CfpVoter.Vote(ctx, cfp.VoteRequest{
		ConfName: input.ConfName,
		TalkID:   int32(input.TalkID),
		GlobalIP: []byte{},
	}); err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return nil, nil
}

// VoteCounts is the resolver for the voteCounts field.
func (r *queryResolver) VoteCounts(ctx context.Context, confName *string) ([]*model.VoteCount, error) {
	ctx = sqlhelper.WithSqlHelper(ctx, r.sqlHelper)

	counts, err := r.CfpVoter.GetCount(ctx, cfp.GetCountRequest{
		ConfName: *confName,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	var resp []*model.VoteCount
	for _, v := range counts {
		resp = append(resp, &model.VoteCount{
			TalkID: int(v.TalkId),
			Count:  v.Count,
		})
	}

	return resp, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
