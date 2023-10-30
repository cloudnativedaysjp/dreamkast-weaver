package repo

import (
	"context"
	"database/sql"
	"dreamkast-weaver/internal/cfp/domain"
	"dreamkast-weaver/internal/cfp/value"
	"dreamkast-weaver/internal/stacktrace"
	"fmt"
	"net"
)

type CfpRepoImpl struct {
	q *Queries
}

func NewCfpRepo(db DBTX) domain.CfpRepo {
	q := New(db)
	return &CfpRepoImpl{q}
}

var _ domain.CfpRepo = (*CfpRepoImpl)(nil)

func (r *CfpRepoImpl) ListCfpVotes(ctx context.Context, confName value.ConfName, vt value.VotingTerm) (*domain.CfpVotes, error) {
	s, e := vt.Value()
	req := ListCfpVotesParams{
		ConferenceName: string(confName.Value()),
		Start:          s,
		End:            e,
	}

	votes, err := r.q.ListCfpVotes(ctx, req)
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("list cfp vote: %w", err))
	}

	return cfpVoteConv.fromDB(votes)
}

func (r *CfpRepoImpl) InsertCfpVote(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error {
	req := InsertCfpVoteParams{
		ConferenceName: string(confName.Value()),
		TalkID:         talkID.Value(),
		ClientIp: sql.NullString{
			String: clientIp.String(),
			Valid:  true,
		},
	}

	if err := r.q.InsertCfpVote(ctx, req); err != nil {
		return stacktrace.With(fmt.Errorf("incert cfp vote: %w", err))
	}
	return nil
}

var cfpVoteConv _cfpVoteConv

type _cfpVoteConv struct{}

func (_cfpVoteConv) fromDB(v []CfpVote) (*domain.CfpVotes, error) {
	conv := func(v *CfpVote) (*domain.CfpVote, error) {
		talkID, err := value.NewTalkID(v.TalkID)
		if err != nil {
			return nil, err
		}
		ip := net.ParseIP(v.ClientIp.String)
		return &domain.CfpVote{
			TalkID:    talkID,
			ClientIp:  ip,
			CreatedAt: v.CreatedAt,
		}, nil
	}

	var items []domain.CfpVote
	for _, p := range v {
		cv := p
		dcv, err := conv(&cv)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view event from DB: %w", err))
		}
		items = append(items, *dcv)
	}

	return &domain.CfpVotes{Items: items}, nil
}
