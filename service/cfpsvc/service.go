package cfpsvc

import (
	"context"

	"github.com/ServiceWeaver/weaver"

	"dreamkast-weaver/service/mysqlsvc"
)

type Cfp struct {
	weaver.AutoMarshal
}

// CfpSvc component.
type T interface {
	Vote(context.Context, string) (string, error)
	Show(context.Context) ([]mysqlsvc.VotingResultItem, error)
}

// Implementation of the CfpSvc component.
type impl struct {
	weaver.Implements[T]
	mySqlSvc mysqlsvc.T
}

func (s *impl) Init(context.Context) error {
	mss, err := weaver.Get[mysqlsvc.T](s)
	s.mySqlSvc = mss
	return err
}

func (s *impl) Vote(ctx context.Context, str string) (string, error) {
	err := s.mySqlSvc.InsertCfpVote(ctx)
	return "ok", err
}

func (s *impl) Show(ctx context.Context) ([]mysqlsvc.VotingResultItem, error) {
	return s.mySqlSvc.ListCfpVotes(ctx)
}
