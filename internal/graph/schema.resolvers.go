package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.29

import (
	"context"
	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/graph/middleware"
	"dreamkast-weaver/internal/graph/model"
	"errors"
	"net"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, input model.VoteInput) (*bool, error) {
	req := cfp.VoteRequest{
		ConfName: input.ConfName.String(),
		TalkID:   input.TalkID,
		GlobalIP: net.ParseIP(middleware.RealIPFromContext(ctx)),
	}

	if err := r.CfpService.Vote(ctx, req); err != nil {
		return nil, err
	}
	return nil, nil
}

// StampOnline is the resolver for the stampOnline field.
func (r *mutationResolver) StampOnline(ctx context.Context, input model.StampOnlineInput) (*bool, error) {
	var e, err error
	profile, err := newProfile(input.ConfName, input.ProfileID)
	err = errors.Join(err, e)

	slotID, e := value.NewSlotID(int32(input.SlotID))
	err = errors.Join(err, e)
	if err != nil {
		return nil, err
	}

	if err := r.DkUiService.StampOnline(ctx, profile, slotID); err != nil {
		return nil, err
	}
	return nil, nil
}

// StampOnSite is the resolver for the stampOnSite field.
func (r *mutationResolver) StampOnSite(ctx context.Context, input model.StampOnSiteInput) (*bool, error) {
	var e, err error
	profile, err := newProfile(input.ConfName, input.ProfileID)
	err = errors.Join(err, e)

	req := dkui.StampRequest{}
	req.TrackID, e = value.NewTrackID(int32(input.TrackID))
	err = errors.Join(err, e)
	req.TalkID, e = value.NewTalkID(int32(input.TalkID))
	err = errors.Join(err, e)
	req.SlotID, e = value.NewSlotID(int32(input.SlotID))
	err = errors.Join(err, e)
	if err != nil {
		return nil, err
	}

	if err := r.DkUiService.StampOnSite(ctx, profile, req); err != nil {
		return nil, err
	}
	return nil, nil
}

// CreateViewEvent is the resolver for the createViewEvent field.
func (r *mutationResolver) CreateViewEvent(ctx context.Context, input model.CreateViewEventInput) (*bool, error) {
	var e, err error
	profile, err := newProfile(input.ConfName, input.ProfileID)
	err = errors.Join(err, e)

	req := dkui.CreateViewEventRequest{}
	req.TrackID, e = value.NewTrackID(int32(input.TrackID))
	err = errors.Join(err, e)
	req.TalkID, e = value.NewTalkID(int32(input.TalkID))
	err = errors.Join(err, e)
	req.SlotID, e = value.NewSlotID(int32(input.SlotID))
	err = errors.Join(err, e)
	if err != nil {
		return nil, err
	}

	if err := r.DkUiService.CreateViewEvent(ctx, profile, req); err != nil {
		return nil, err
	}
	return nil, nil
}

// VoteCounts is the resolver for the voteCounts field.
func (r *queryResolver) VoteCounts(ctx context.Context, confName model.ConfName) ([]*model.VoteCount, error) {
	resp, err := r.CfpService.VoteCounts(ctx, confName.String())
	if err != nil {
		return nil, err
	}

	var counts []*model.VoteCount
	for _, v := range resp {
		counts = append(counts, &model.VoteCount{
			TalkID: v.TalkID,
			Count:  v.Count,
		})
	}

	return counts, nil
}

// ViewingSlots is the resolver for the viewingSlots field.
func (r *queryResolver) ViewingSlots(ctx context.Context, confName model.ConfName, profileID int) ([]*model.ViewingSlot, error) {
	profile, err := newProfile(confName, profileID)
	if err != nil {
		return nil, err
	}

	devents, err := r.DkUiService.ViewingEvents(ctx, profile)
	if err != nil {
		return nil, err
	}

	var viewingSlots []*model.ViewingSlot
	for k, v := range devents.ViewingSeconds() {
		viewingSlots = append(viewingSlots, &model.ViewingSlot{
			SlotID:      int(k.Value()),
			ViewingTime: int(v),
		})
	}

	return viewingSlots, nil
}

// StampChallenges is the resolver for the stampChallenges field.
func (r *queryResolver) StampChallenges(ctx context.Context, confName model.ConfName, profileID int) ([]*model.StampChallenge, error) {
	profile, err := newProfile(confName, profileID)
	if err != nil {
		return nil, err
	}

	dstamps, err := r.DkUiService.StampChallenges(ctx, profile)
	if err != nil {
		return nil, err
	}

	var stamps []*model.StampChallenge
	for _, dst := range dstamps.Items {
		stamps = append(stamps, &model.StampChallenge{
			SlotID:    int(dst.SlotID.Value()),
			Condition: model.ChallengeCondition(dst.Condition.Value()),
			UpdatedAt: int(dst.UpdatedAt.Unix()),
		})
	}

	return stamps, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func newProfile(confName model.ConfName, profileID int) (dkui.Profile, error) {
	var e, err error
	profile := dkui.Profile{}
	profile.ConfName, e = value.NewConfName(value.ConferenceKind(confName))
	err = errors.Join(err, e)
	profile.ID, e = value.NewProfileID(int32(profileID))
	err = errors.Join(err, e)

	return profile, err
}
