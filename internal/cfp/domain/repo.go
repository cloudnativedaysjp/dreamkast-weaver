package domain

import (
	"context"
	"dreamkast-weaver/internal/cfp/value"
	"net"
)

type CfpRepo interface {
	ListCfpVotes(ctx context.Context, confName value.ConfName, vt *VotingTerm) (*CfpVotes, error)
	InsertCfpVote(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error
}
