package domain

import (
	"context"
	"net"
	"time"

	"dreamkast-weaver/internal/dkui/value"
)

type DkUiRepo interface {
	ListViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*ViewEvents, error)
	InsertViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *ViewEvent) error

	GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*StampChallenges, error)
	UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *StampChallenges) error

	InsertTrackViewer(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) (err error)
	ListTrackViewer(ctx context.Context, from, to time.Time) (*TrackViewers, error)
}

type CfpRepo interface {
	ListCfpVotes(ctx context.Context, confName value.ConfName, vt value.VotingTerm) (*CfpVotes, error)
	InsertCfpVote(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error
}
