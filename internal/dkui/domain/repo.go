package domain

import (
	"context"
	"time"

	"dreamkast-weaver/internal/dkui/value"
)

type DkUiRepo interface {
	ListViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*ViewEvents, error)
	InsertViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *ViewEvent) error

	GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*StampChallenges, error)
	UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *StampChallenges) error

	UpsertViewerCount(ctx context.Context, cn value.ConfName, vc ViewerCount) error
	ListViewerCounts(ctx context.Context, cn value.ConfName) (*ViewerCounts, error)
	InsertTrackViewer(ctx context.Context, profileID value.ProfileID, trackName value.TrackName) (err error)
	ListTrackViewer(ctx context.Context, from, to time.Time) (*TrackViewers, error)
}
