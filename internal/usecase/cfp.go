package cfp

import (
	"context"
	"net"

	"dreamkast-weaver/internal/cfp/domain"
	"dreamkast-weaver/internal/cfp/repo"
	"dreamkast-weaver/internal/cfp/value"
	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/stacktrace"

	"github.com/ServiceWeaver/weaver"
)

type Service interface {
	Vote(ctx context.Context, req VoteRequest) error
	VoteCounts(ctx context.Context, req VoteCountsRequest) ([]*domain.VoteCount, error)
}

type VoteRequest struct {
	weaver.AutoMarshal
	ConfName value.ConfName
	TalkID   value.TalkID
	ClientIp net.IP
}

type VoteCountsRequest struct {
	weaver.AutoMarshal
	ConfName    value.ConfName
	VotingTerm  value.VotingTerm
	SpanSeconds value.SpanSeconds
}

type ServiceImpl struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]

	sh     *sqlhelper.SqlHelper
	domain domain.CfpDomain
}

var _ Service = (*ServiceImpl)(nil)

type config struct {
	DBUser     string `toml:"db_user"`
	DBPassword string `toml:"db_password"`
	DBEndpoint string `toml:"db_endpoint"`
	DBPort     string `toml:"db_port"`
	DBName     string `toml:"db_name"`
}

func (c *config) SqlOption() *sqlhelper.SqlOption {
	return &sqlhelper.SqlOption{
		User:     c.DBUser,
		Password: c.DBPassword,
		Endpoint: c.DBEndpoint,
		Port:     c.DBPort,
		DbName:   c.DBName,
	}
}

func NewService(sh *sqlhelper.SqlHelper) Service {
	return &ServiceImpl{sh: sh}
}

func (s *ServiceImpl) Init(ctx context.Context) error {
	opt := s.Config().SqlOption()
	if err := opt.Validate(); err != nil {
		opt = sqlhelper.NewOptionFromEnv("cfp")
	}
	sh, err := sqlhelper.NewSqlHelper(opt)
	if err != nil {
		return err
	}
	s.sh = sh
	return nil
}

func (s *ServiceImpl) HandleError(msg string, err error) {
	if err != nil {
		if derrors.IsUserError(err) {
			s.Logger().With("errorType", "user-side").Info(msg, err)
		} else {
			s.Logger().With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *ServiceImpl) VoteCounts(ctx context.Context, req VoteCountsRequest) (resp []*domain.VoteCount, err error) {
	defer func() {
		s.HandleError("get voteCounts", err)
	}()

	r := repo.NewCfpRepo(s.sh.DB())

	dvotes, err := r.ListCfpVotes(ctx, req.ConfName, req.VotingTerm)
	if err != nil {
		return nil, err
	}

	dvc := s.domain.TallyCfpVotes(dvotes, req.SpanSeconds)

	return dvc, nil
}

func (s *ServiceImpl) Vote(ctx context.Context, req VoteRequest) (err error) {
	defer func() {
		s.HandleError("vote", err)
	}()

	r := repo.NewCfpRepo(s.sh.DB())

	return r.InsertCfpVote(ctx, req.ConfName, req.TalkID, req.ClientIp)
}
