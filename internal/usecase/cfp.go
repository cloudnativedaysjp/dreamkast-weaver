package usecase

import (
	"context"
	"net"

	derrors "dreamkast-weaver/internal/domain/errors"
	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/repo"
	"dreamkast-weaver/internal/logger"
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/stacktrace"
)

type CfpService interface {
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

type CfpServiceImpl struct {
	sh        *sqlhelper.SqlHelper
	cfpDomain dmodel.CfpDomain
}

var _ CfpService = (*CfpServiceImpl)(nil)

func NewCFPService(sh *sqlhelper.SqlHelper) CfpService {
	return &CfpServiceImpl{sh: sh}
}

func (s *CfpServiceImpl) Init(ctx context.Context) error {
	return nil
}

func (s *CfpServiceImpl) HandleError(ctx context.Context, msg string, err error) {
	logger := logger.FromCtx(ctx)
	if err != nil {
		if derrors.IsUserError(err) {
			logger.With("errorType", "user-side").Info(msg, err)
		} else {
			logger.With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *CfpServiceImpl) VoteCounts(ctx context.Context, req VoteCountsRequest) (resp []*dmodel.VoteCount, err error) {
	defer func() {
		s.HandleError(ctx, "get voteCounts", err)
	}()

	r := repo.NewCfpRepo(s.sh.DB())

	dvotes, err := r.ListCfpVotes(ctx, req.ConfName, req.VotingTerm)
	if err != nil {
		return nil, err
	}

	dvc := s.cfpDomain.TallyCfpVotes(dvotes, req.SpanSeconds)

	return dvc, nil
}

func (s *CfpServiceImpl) Vote(ctx context.Context, req VoteRequest) (err error) {
	defer func() {
		s.HandleError(ctx, "vote", err)
	}()

	r := repo.NewCfpRepo(s.sh.DB())

	return r.InsertCfpVote(ctx, req.ConfName, req.TalkID, req.ClientIp)
}
