package cfp

import (
	"context"
	"database/sql"
	"fmt"

	"dreamkast-weaver/internal/cfp/repo"
	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/graph/model"
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/stacktrace"

	"github.com/ServiceWeaver/weaver"
)

//go:generate go run github.com/ServiceWeaver/weaver/cmd/weaver generate .

type Service interface {
	Vote(ctx context.Context, input model.VoteInput) error
	VoteCounts(ctx context.Context, confName model.ConfName) ([]*model.VoteCount, error)
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
	if err != nil && !derrors.IsUserError(err) {
		s.Logger().With("stacktrace", stacktrace.Get(err)).Error(msg, err)
	}
}

func (s *ServiceImpl) VoteCounts(ctx context.Context, confName model.ConfName) (resp []*model.VoteCount, err error) {
	defer func() {
		s.HandleError("get voteCounts", err)
	}()

	r := repo.New(s.sh.DB())

	votes, err := r.ListCfpVotes(ctx, confName.String())
	if err != nil {
		return nil, fmt.Errorf("list cfp vote: %w", err)
	}

	// TODO move to domain package
	counts := map[int32]int{}
	for _, vote := range votes {
		counts[vote.TalkID]++
	}

	for talkID, count := range counts {
		resp = append(resp, &model.VoteCount{
			TalkID: int(talkID),
			Count:  count,
		})
	}

	return resp, nil
}

func (s *ServiceImpl) Vote(ctx context.Context, input model.VoteInput) (err error) {
	defer func() {
		s.HandleError("vote", err)
	}()

	r := repo.New(s.sh.DB())

	if err := r.InsertCfpVote(ctx, repo.InsertCfpVoteParams{
		ConferenceName: input.ConfName.String(),
		TalkID:         int32(input.TalkID),
		GlobalIp: sql.NullString{
			String: input.GlobalIP,
			Valid:  true,
		},
	}); err != nil {
		return fmt.Errorf("incert cfp vote: %w", err)
	}

	return nil
}
