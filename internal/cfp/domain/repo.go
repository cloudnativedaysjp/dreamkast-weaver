package domain

import (
	"context"
	"dreamkast-weaver/internal/dkui/value"
)

type CfpRepo interface {
	ListCfpVotes(ctx context.Context, confName value.ConfName) (*CfpVotes, error)
	InsertCfpVote(ctx context.Context, cfpVote *CfpVote) error
}
