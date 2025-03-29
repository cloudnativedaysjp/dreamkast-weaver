package application

import (
	"context"
	"net"

	derrors "dreamkast-weaver/internal/domain/errors"
	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/repo"
	"dreamkast-weaver/internal/pkg/logger"
	"dreamkast-weaver/internal/pkg/sqlhelper"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type CfpApp interface {
	Vote(ctx context.Context, req VoteRequest) error
	VoteCounts(ctx context.Context, req VoteCountsRequest) ([]*dmodel.VoteCount, error)
}

type VoteRequest struct {
	ConfName value.ConfName
	TalkID   value.TalkID
	ClientIp net.IP
}

type VoteCountsRequest struct {
	ConfName    value.ConfName
	VotingTerm  value.VotingTerm
	SpanSeconds value.SpanSeconds
}

type CfpAppImpl struct {
	sh *sqlhelper.SqlHelper
}

var _ CfpApp = (*CfpAppImpl)(nil)

func NewCfpApp(sh *sqlhelper.SqlHelper) CfpApp {
	return &CfpAppImpl{sh: sh}
}

func (s *CfpAppImpl) handleError(ctx context.Context, msg string, err error) {
	logger := logger.FromCtx(ctx)
	if err != nil {
		if derrors.IsUserError(err) {
			logger.With("errorType", "user-side").Info(msg, err)
		} else {
			logger.With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *CfpAppImpl) VoteCounts(ctx context.Context, req VoteCountsRequest) (resp []*dmodel.VoteCount, err error) {
	defer func() {
		s.handleError(ctx, "get voteCounts", err)
	}()

	r := repo.NewCfpVoteRepo(s.sh.DB())

	dvotes, err := r.List(ctx, req.ConfName, req.VotingTerm)
	if err != nil {
		return nil, err
	}

	dvc := dvotes.Tally(req.SpanSeconds)

	return dvc, nil
}

func (s *CfpAppImpl) Vote(ctx context.Context, req VoteRequest) (err error) {
	defer func() {
		s.handleError(ctx, "vote", err)
	}()

	r := repo.NewCfpVoteRepo(s.sh.DB())

	return r.Insert(ctx, req.ConfName, req.TalkID, req.ClientIp)
}
