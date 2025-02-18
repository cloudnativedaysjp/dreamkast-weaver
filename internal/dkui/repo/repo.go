package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/stacktrace"
)

type DkUiRepoImpl struct {
	q *Queries
}

func NewDkUiRepo(db DBTX) domain.DkUiRepo {
	q := New(db)
	return &DkUiRepoImpl{q}
}

var _ domain.DkUiRepo = (*DkUiRepoImpl)(nil)

func (r *DkUiRepoImpl) GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*domain.StampChallenges, error) {
	data, err := r.q.GetTrailmapStamps(ctx, GetTrailmapStampsParams{
		ConferenceName: confName.String(),
		ProfileID:      profileID.Value(),
	})
	if err != nil {
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, stacktrace.With(fmt.Errorf("get stamp challenges: %w", err))
			}
		}
	}

	return stampChallengeConv.fromDB(data.Stamps)
}

func (r *DkUiRepoImpl) InsertViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *domain.ViewEvent) error {
	if err := r.q.InsertViewEvents(ctx, InsertViewEventsParams{
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

func (r *DkUiRepoImpl) ListViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*domain.ViewEvents, error) {
	data, err := r.q.ListViewEvents(ctx, ListViewEventsParams{
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

func (r *DkUiRepoImpl) UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *domain.StampChallenges) error {
	buf, err := stampChallengeConv.toDB(scs)
	if err != nil {
		return stacktrace.With(err)
	}

	if err := r.q.UpsertTrailmapStamp(ctx, UpsertTrailmapStampParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
		Stamps:         buf,
	}); err != nil {
		return stacktrace.With(fmt.Errorf("upsert stamp challenges: %w", err))
	}
	return nil
}

var stampChallengeConv _stampChallengeConv

type _stampChallengeConv struct{}

type _stampChallenge struct {
	SlotID    int32
	Condition string
	UpdatedAt time.Time
}

func (_stampChallengeConv) toDB(v *domain.StampChallenges) (json.RawMessage, error) {
	conv := func(dsc *domain.StampChallenge) *_stampChallenge {
		return &_stampChallenge{
			SlotID:    dsc.SlotID.Value(),
			Condition: string(dsc.Condition.Value()),
			UpdatedAt: dsc.UpdatedAt,
		}
	}

	var stamps []_stampChallenge
	for _, p := range v.Items {
		st := p
		stamps = append(stamps, *conv(&st))
	}

	buf, err := json.Marshal(stamps)
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("convert stamp challenges to DB: %w", err))
	}
	return json.RawMessage(buf), nil
}

func (_stampChallengeConv) fromDB(v json.RawMessage) (*domain.StampChallenges, error) {
	conv := func(sc *_stampChallenge) (*domain.StampChallenge, error) {
		slotID, err := value.NewSlotID(sc.SlotID)
		if err != nil {
			return nil, err
		}
		cond, err := value.NewStampCondition(value.StampConditionKind(sc.Condition))
		if err != nil {
			return nil, err
		}

		return &domain.StampChallenge{
			SlotID:    slotID,
			Condition: cond,
			UpdatedAt: sc.UpdatedAt,
		}, nil
	}

	if v == nil {
		return &domain.StampChallenges{}, nil
	}

	var stamps []_stampChallenge
	if err := json.Unmarshal(v, &stamps); err != nil {
		return nil, stacktrace.With(fmt.Errorf("convert stamp challenges from DB: %w", err))
	}

	var items []domain.StampChallenge
	for _, p := range stamps {
		st := p
		dst, err := conv(&st)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert stamp challenges from DB: %w", err))
		}
		items = append(items, *dst)
	}

	return &domain.StampChallenges{Items: items}, nil
}

var viewEventConv _viewEventConv

type _viewEventConv struct{}

func (_viewEventConv) fromDB(v []ViewEvent) (*domain.ViewEvents, error) {
	conv := func(v *ViewEvent) (*domain.ViewEvent, error) {
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
		return &domain.ViewEvent{
			TrackID:        trackID,
			TalkID:         talkID,
			SlotID:         slotID,
			ViewingSeconds: viewingSeconds,
			CreatedAt:      v.CreatedAt,
		}, nil
	}

	var items []domain.ViewEvent

	for _, p := range v {
		ev := p
		dev, err := conv(&ev)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view event from DB: %w", err))
		}
		items = append(items, *dev)
	}

	return &domain.ViewEvents{Items: items}, nil
}

// func (_viewEventConv) toDB(confName value.ConfName, profileID value.ProfileID, v *domain.ViewEvents) []ViewEvent {
// 	conv := func(dev *domain.ViewEvent) *ViewEvent {
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

func (r *DkUiRepoImpl) InsertTrackViewer(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) error {
	if err := r.q.InsertTrackViewer(ctx, InsertTrackViewerParams{
		ProfileID: profileID.Value(),
		TrackName: trackName.String(),
		TalkID:    talkID.Value(),
	}); err != nil {
		return stacktrace.With(fmt.Errorf("insert viewing track: %w", err))
	}
	return nil
}

func (r *DkUiRepoImpl) ListTrackViewer(ctx context.Context, from, to time.Time) (*domain.TrackViewers, error) {
	tvs, err := r.q.ListTrackViewer(ctx, ListTrackViewerParams{
		FromCreatedAt: from,
		ToCreatedAt:   to,
	})
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("list viewing track: %w", err))
	}

	return trackViewerConv.fromDB(tvs)
}

var trackViewerConv _trackViewerConv

type _trackViewerConv struct{}

func (_trackViewerConv) fromDB(tvs []TrackViewer) (*domain.TrackViewers, error) {
	conv := func(v *TrackViewer) (*domain.TrackViewer, error) {
		tn, err := value.NewTrackName(v.TrackName)
		if err != nil {
			return nil, err
		}
		pID, err := value.NewProfileID(v.ProfileID)
		if err != nil {
			return nil, err
		}

		return &domain.TrackViewer{
			TrackName: tn,
			ProfileID: pID,
			CreatedAt: v.CreatedAt,
		}, nil
	}

	var items []domain.TrackViewer

	for _, tv := range tvs {
		v := tv
		dv, err := conv(&v)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert track viewer from DB: %w", err))
		}
		items = append(items, *dv)
	}

	return &domain.TrackViewers{Items: items}, nil
}

// func (_trackViewerConv) toDB(v *domain.TrackViewers) []TrackViewer {
// 	conv := func(dtv *domain.TrackViewer) *TrackViewer {
// 		return &TrackViewer{
// 			TrackName: dtv.TrackName.String(),
// 			ProfileID: dtv.ProfileID.Value(),
// 			CreatedAt: dtv.CreatedAt,
// 		}
// 	}
// 	var tvs []TrackViewer
// 	for _, vc := range v.Items {
// 		tvs = append(tvs, *conv(&vc))
// 	}
// 	return tvs
// }
