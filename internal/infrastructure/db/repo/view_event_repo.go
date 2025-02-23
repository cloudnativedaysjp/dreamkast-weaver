package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/repo"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/dbgen"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type ViewEventRepoImpl struct {
	q *dbgen.Queries
}

func NewViewEventRepo(db dbgen.DBTX) repo.ViewEventRepo {
	q := dbgen.New(db)
	return &ViewEventRepoImpl{q}
}

var _ repo.ViewEventRepo = (*ViewEventRepoImpl)(nil)

func (r *ViewEventRepoImpl) Insert(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *dmodel.ViewEvent) error {
	if err := r.q.InsertViewEvents(ctx, dbgen.InsertViewEventsParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
		TrackID:        ev.TrackID.Value(),
		TalkID:         ev.TalkID.Value(),
		SlotID:         ev.SlotID.Value(),
		ViewingSeconds: ev.ViewingSeconds.Value(),
	}); err != nil {
		return stacktrace.With(fmt.Errorf("insert view event: %w", err))
	}
	return nil
}

func (r *ViewEventRepoImpl) List(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*dmodel.ViewEvents, error) {
	data, err := r.q.ListViewEvents(ctx, dbgen.ListViewEventsParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, stacktrace.With(fmt.Errorf("list view event: %w", err))
		}
	}

	return viewEventConv.fromDB(data)
}

var viewEventConv _viewEventConv

type _viewEventConv struct{}

func (_viewEventConv) fromDB(v []dbgen.ViewEvent) (*dmodel.ViewEvents, error) {
	conv := func(v *dbgen.ViewEvent) (*dmodel.ViewEvent, error) {
		trackID, err := value.NewTrackID(v.TrackID)
		if err != nil {
			return nil, err
		}
		talkID, err := value.NewTalkID(v.TalkID)
		if err != nil {
			return nil, err
		}
		slotID, err := value.NewSlotID(v.SlotID)
		if err != nil {
			return nil, err
		}
		viewingSeconds, err := value.NewViewingSeconds(v.ViewingSeconds)
		if err != nil {
			return nil, err
		}
		return &dmodel.ViewEvent{
			TrackID:        trackID,
			TalkID:         talkID,
			SlotID:         slotID,
			ViewingSeconds: viewingSeconds,
			CreatedAt:      v.CreatedAt,
		}, nil
	}

	var items []dmodel.ViewEvent

	for _, p := range v {
		ev := p
		dev, err := conv(&ev)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view event from DB: %w", err))
		}
		items = append(items, *dev)
	}

	return &dmodel.ViewEvents{Items: items}, nil
}

// func (_viewEventConv) toDB(confName value.ConfName, profileID value.ProfileID, v *dmodel.ViewEvents) []ViewEvent {
// 	conv := func(dev *dmodel.ViewEvent) *ViewEvent {
// 		return &ViewEvent{
// 			ConferenceName: string(confName.Value()),
// 			ProfileID:      profileID.Value(),
// 			TrackID:        dev.TrackID.Value(),
// 			TalkID:         dev.TalkID.Value(),
// 			SlotID:         dev.SlotID.Value(),
// 			ViewingSeconds: dev.ViewingSeconds.Value(),
// 			CreatedAt:      dev.CreatedAt,
// 		}
// 	}
//
// 	var events []ViewEvent
// 	for _, ev := range v.Items {
// 		events = append(events, *conv(&ev))
// 	}
//
// 	return events
// }
