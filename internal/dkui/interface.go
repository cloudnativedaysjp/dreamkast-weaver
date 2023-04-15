package dkui

import (
	"context"
	"dreamkast-weaver/internal/graph/model"

	"github.com/ServiceWeaver/weaver"
)

type Service interface {
	CreateWatchEvent(ctx context.Context, req CreateWatchEventInput) error
	StampOnline(ctx context.Context, req StampOnlineInput) error
	StampOnSite(ctx context.Context, req StampOnSiteInput) error
	ViewingSlots(ctx context.Context, confName model.ConfName, profileID int) ([]*ViewingSlot, error)
	StampChallenges(ctx context.Context, confName model.ConfName, profileID int) ([]*StampChallenge, error)
}

type CreateWatchEventInput struct {
	weaver.AutoMarshal
	model.CreateWatchEventInput
}

func NewCreateWatchEventInput(v model.CreateWatchEventInput) *CreateWatchEventInput {
	return &CreateWatchEventInput{
		CreateWatchEventInput: v,
	}
}

type StampOnlineInput struct {
	weaver.AutoMarshal
	model.StampOnlineInput
}

func NewStampOnlineInput(v model.StampOnlineInput) *StampOnlineInput {
	return &StampOnlineInput{
		StampOnlineInput: v,
	}
}

type StampOnSiteInput struct {
	weaver.AutoMarshal
	model.StampOnSiteInput
}

func NewStampOnSiteInput(v model.StampOnSiteInput) *StampOnSiteInput {
	return &StampOnSiteInput{
		StampOnSiteInput: v,
	}
}

type ViewingSlot struct {
	weaver.AutoMarshal
	model.ViewingSlot
}

func NewViewingSlot(v model.ViewingSlot) *ViewingSlot {
	return &ViewingSlot{
		ViewingSlot: v,
	}
}

type StampChallenge struct {
	weaver.AutoMarshal
	model.StampChallenge
}

func NewStampChallenge(v model.StampChallenge) *StampChallenge {
	return &StampChallenge{
		StampChallenge: v,
	}
}
