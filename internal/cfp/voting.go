package cfp

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"dreamkast-weaver/internal/cfp/repo"
	"dreamkast-weaver/internal/sqlhelper"

	"github.com/ServiceWeaver/weaver"
)

type Voter interface {
	Vote(ctx context.Context, req VoteRequest) error
	GetCount(ctx context.Context, req GetCountRequest) (GetCountResponse, error)
}

// TODO validate
type VoteRequest struct {
	weaver.AutoMarshal
	ConfName string
	TalkID   int32
	GlobalIP net.IP
}

type GetCountRequest struct {
	weaver.AutoMarshal
	ConfName string
}

type VoteCount struct {
	weaver.AutoMarshal
	TalkID int32
	Count  int
}

type GetCountResponse []VoteCount

// VoterImpl implements cfp.Voter
type VoterImpl struct {
	weaver.Implements[Voter]
	weaver.WithConfig[config]

	sh *sqlhelper.SqlHelper
}

var _ Voter = (*VoterImpl)(nil)

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

func NewVoter(sh *sqlhelper.SqlHelper) Voter {
	return &VoterImpl{sh: sh}
}

func (v *VoterImpl) Init(ctx context.Context) error {
	cfg := v.Config()
	log.Printf("config: %#v\n", cfg)
	sh, err := sqlhelper.NewSqlHelper(cfg.SqlOption())
	if err != nil {
		return err
	}
	v.sh = sh
	return nil
}

func (v *VoterImpl) GetCount(ctx context.Context, req GetCountRequest) (GetCountResponse, error) {
	r := repo.New(v.sh.DB())

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
		resp = append(resp, VoteCount{
			TalkID: talkID,
			Count:  count,
		})
	}

	return resp, nil
}

func (v *VoterImpl) Vote(ctx context.Context, req VoteRequest) error {
	r := repo.New(v.sh.DB())

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
