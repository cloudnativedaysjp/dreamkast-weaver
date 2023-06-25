package domain

import (
	"context"

	"dreamkast-weaver/internal/dkui/value"
)

type DkUiRepo interface {
	ListViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*ViewEvents, error)
	InsertViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *ViewEvent) error

	GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*StampChallenges, error)
	UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *StampChallenges) error

	UpsertViewerCount(ctx context.Context, cn value.ConfName, vc ViewerCount) error
	GetViewerCount(ctx context.Context, cn value.ConfName, trackID value.TrackID) (*ViewerCount, error)
}
