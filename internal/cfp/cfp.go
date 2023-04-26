package cfp

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"dreamkast-weaver/internal/cfp/repo"
	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/stacktrace"

	"github.com/ServiceWeaver/weaver"
)

type Service interface {
	Vote(ctx context.Context, req VoteRequest) error
	VoteCounts(ctx context.Context, confName string) ([]*VoteCount, error)
}

type VoteRequest struct {
	weaver.AutoMarshal
	ConfName string
	TalkID   int
	ClientIP net.IP
}

type VoteCount struct {
	weaver.AutoMarshal
	TalkID int
	Count  int
}

type ServiceImpl struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]

	sh *sqlhelper.SqlHelper
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

func (s *ServiceImpl) VoteCounts(ctx context.Context, confName string) (resp []*VoteCount, err error) {
	defer func() {
		s.HandleError("get voteCounts", err)
	}()

	r := repo.New(s.sh.DB())

	votes, err := r.ListCfpVotes(ctx, confName)
	if err != nil {
		return nil, fmt.Errorf("list cfp vote: %w", err)
	}

	// TODO move to domain package
	counts := map[int32]int{}
	for _, vote := range votes {
		counts[vote.TalkID]++
	}

	for talkID, count := range counts {
		resp = append(resp, &VoteCount{
			TalkID: int(talkID),
			Count:  count,
		})
	}

	return resp, nil
}

func (s *ServiceImpl) Vote(ctx context.Context, req VoteRequest) (err error) {
	defer func() {
		s.HandleError("vote", err)
	}()

	r := repo.New(s.sh.DB())

	if err := r.InsertCfpVote(ctx, repo.InsertCfpVoteParams{
		ConferenceName: req.ConfName,
		TalkID:         int32(req.TalkID),
		ClientIp: sql.NullString{
			String: req.ClientIP.String(),
			Valid:  true,
		},
	}); err != nil {
		return fmt.Errorf("incert cfp vote: %w", err)
	}

	return nil
}
