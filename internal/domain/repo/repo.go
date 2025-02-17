package repo

import (
	"context"
	"net"
	"time"

	"dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
)

type DkUiRepo interface {
	ListViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*model.ViewEvents, error)
	InsertViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *model.ViewEvent) error

	GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*model.StampChallenges, error)
	UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *model.StampChallenges) error

	InsertTrackViewer(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) (err error)
	ListTrackViewer(ctx context.Context, from, to time.Time) (*model.TrackViewers, error)
}

type CfpRepo interface {
	ListCfpVotes(ctx context.Context, confName value.ConfName, vt value.VotingTerm) (*model.CfpVotes, error)
	InsertCfpVote(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error
}
