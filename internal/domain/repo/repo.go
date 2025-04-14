package repo

import (
	"context"
	"net"
	"time"

	"dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
)

type CfpVoteRepo interface {
	List(ctx context.Context, confName value.ConfName, vt value.VotingTerm) (*model.CfpVotes, error)
	Insert(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error
}

type TrackViewerRepo interface {
	Insert(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) (err error)
	List(ctx context.Context, from, to time.Time) (*model.TrackViewers, error)
}

type TrailMapStampsRepo interface {
	Get(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*model.StampChallenges, error)
	Upsert(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *model.StampChallenges) error
}

type ViewEventRepo interface {
	List(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*model.ViewEvents, error)
	Insert(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *model.ViewEvent) error
}
