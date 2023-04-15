package dkui

import (
	"context"
	"dreamkast-weaver/internal/graph/model"
)

type Service interface {
	CreateWatchEvent(ctx context.Context, req model.CreateWatchEventInput) error
	StampOnline(ctx context.Context, req model.StampOnlineInput) error
	StampOnSite(ctx context.Context, req model.StampOnSiteInput) error
	ViewingSlots(ctx context.Context, confName model.ConfName, profileID int) ([]*model.ViewingSlot, error)
	StampChallenges(ctx context.Context, confName model.ConfName, profileID int) ([]*model.StampChallenge, error)
}
