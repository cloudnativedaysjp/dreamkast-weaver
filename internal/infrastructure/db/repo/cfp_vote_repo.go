package repo

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/repo"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/dbgen"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type CfpRepoImpl struct {
	q *dbgen.Queries
}

func NewCfpVoteRepo(db dbgen.DBTX) repo.CfpVoteRepo {
	q := dbgen.New(db)
	return &CfpRepoImpl{q}
}

func (r *CfpRepoImpl) List(ctx context.Context, confName value.ConfName, vt value.VotingTerm) (*dmodel.CfpVotes, error) {
	s, e := vt.Value()
	req := dbgen.ListCfpVotesParams{
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

func (r *CfpRepoImpl) Insert(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error {
	req := dbgen.InsertCfpVoteParams{
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

func (_cfpVoteConv) fromDB(v []dbgen.CfpVote) (*dmodel.CfpVotes, error) {
	conv := func(v *dbgen.CfpVote) (*dmodel.CfpVote, error) {
		talkID, err := value.NewTalkID(v.TalkID)
		if err != nil {
			return nil, err
		}
		ip := net.ParseIP(v.ClientIp.String)
		return &dmodel.CfpVote{
			TalkID:    talkID,
			ClientIp:  ip,
			CreatedAt: v.CreatedAt,
		}, nil
	}

	var items []dmodel.CfpVote
	for _, p := range v {
		cv := p
		dcv, err := conv(&cv)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view event from DB: %w", err))
		}
		items = append(items, *dcv)
	}

	return &dmodel.CfpVotes{Items: items}, nil
}
