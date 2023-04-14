package domain

import (
	"context"
	"dreamkast-weaver/internal/dkui/value"
)

type DkUiRepo interface {
	ListWatchEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*WatchEvents, error)
	InsertWatchEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *WatchEvent) error

	GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*StampChallenges, error)
	UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *StampChallenges) error
}
