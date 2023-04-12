package cfp

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"dreamkast-weaver/internal/cfp/repo"
	"dreamkast-weaver/internal/sqlhelper"
)

type Voter interface {
	Vote(ctx context.Context, req VoteRequest) error
	GetCount(ctx context.Context, req GetCountRequest) (GetCountResponse, error)
}

// TODO validate
type VoteRequest struct {
	ConfName string
	TalkID   int32
	GlobalIP net.IP
}

type GetCountRequest struct {
	ConfName string
}

type VoteCount struct {
	TalkId int32
	Count  int
}

type GetCountResponse []VoteCount

// VoterImpl implements cfp.Voter
type VoterImpl struct{}

var _ Voter = (*VoterImpl)(nil)

func NewVoter() Voter {
	return &VoterImpl{}
}

func (*VoterImpl) GetCount(ctx context.Context, req GetCountRequest) (GetCountResponse, error) {
	sqlh := sqlhelper.SqlHelperFromContext(ctx)
	r := repo.New(sqlh.DB())

	votes, err := r.ListCfpVotes(ctx, req.ConfName)
	if err != nil {
		return nil, fmt.Errorf("list cfp vote: %w", err)
	}

	// TODO move to domain package
	counts := map[int32]int{}
	for _, vote := range votes {
		counts[vote.TalkID]++
	}

	var resp GetCountResponse
	for talkID, count := range counts {
		resp = append(resp, VoteCount{talkID, count})
	}

	return resp, nil
}

func (*VoterImpl) Vote(ctx context.Context, req VoteRequest) error {
	sqlh := sqlhelper.SqlHelperFromContext(ctx)
	r := repo.New(sqlh.DB())

	if err := r.InsertCfpVote(ctx, repo.InsertCfpVoteParams{
		ConferenceName: req.ConfName,
		TalkID:         req.TalkID,
		GlobalIp: sql.NullString{
			String: req.GlobalIP.String(),
			Valid:  true,
		},
	}); err != nil {
		return fmt.Errorf("incert cfp vote: %w", err)
	}

	return nil
}
