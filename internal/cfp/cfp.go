package cfp

import (
	"context"
	"database/sql"
	"fmt"
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
	VoteCounts(ctx context.Context, confName value.ConfName) ([]*domain.VoteCount, error)
}

type VoteRequest struct {
	weaver.AutoMarshal
	ConfName value.ConfName
	TalkID   value.TalkID
	ClientIp net.IP
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

func (s *ServiceImpl) VoteCounts(ctx context.Context, confName value.ConfName) (resp []*domain.VoteCount, err error) {
	defer func() {
		s.HandleError("get voteCounts", err)
	}()

	r := repo.New(s.sh.DB())

	votes, err := r.ListCfpVotes(ctx, string(confName.Value()))
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("list cfp vote: %w", err))
	}

	dvotes, err := cfpVoteConv.fromDB(votes)
	if err != nil {
		return nil, err
	}

	dvc := s.domain.TallyCfpVotes(dvotes)

	return dvc, nil
}

func (s *ServiceImpl) Vote(ctx context.Context, req VoteRequest) (err error) {
	defer func() {
		s.HandleError("vote", err)
	}()

	r := repo.New(s.sh.DB())

	if err := r.InsertCfpVote(ctx, repo.InsertCfpVoteParams{
		ConferenceName: string(req.ConfName.Value()),
		TalkID:         req.TalkID.Value(),
		ClientIp: sql.NullString{
			String: req.ClientIp.String(),
			Valid:  true,
		},
	}); err != nil {
		return stacktrace.With(fmt.Errorf("incert cfp vote: %w", err))
	}

	return nil
}

var cfpVoteConv _cfpVoteConv

type _cfpVoteConv struct{}

func (_cfpVoteConv) fromDB(v []repo.CfpVote) (*domain.CfpVotes, error) {
	conv := func(v *repo.CfpVote) (*domain.CfpVote, error) {
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
