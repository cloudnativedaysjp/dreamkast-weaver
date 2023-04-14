package repo

import (
	"context"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"encoding/json"
	"time"
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
		return nil, err
	}

	return stampChallengeConv.fromDB(data.Stamps)
}

func (r *DkUiRepoImpl) InsertWatchEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *domain.WatchEvent) error {
	if err := r.q.InsertWatchEvents(ctx, InsertWatchEventsParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
		TrackID:        ev.TrackID.Value(),
		TalkID:         ev.TalkID.Value(),
		SlotID:         ev.SlotID.Value(),
		ViewingSeconds: ev.ViewingSeconds.Value(),
	}); err != nil {
		return err
	}
	return nil
}

func (r *DkUiRepoImpl) ListWatchEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*domain.WatchEvents, error) {
	data, err := r.q.ListWatchEvents(ctx, ListWatchEventsParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
	})
	if err != nil {
		return nil, err
	}

	return watchEventConv.fromDB(data)
}

func (r *DkUiRepoImpl) UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *domain.StampChallenges) error {
	buf, err := stampChallengeConv.toDB(scs)
	if err != nil {
		return err
	}

	if err := r.q.UpsertTrailmapStamp(ctx, UpsertTrailmapStampParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
		Stamps:         buf,
	}); err != nil {
		return err
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
	for _, st := range v.Items {
		stamps = append(stamps, *conv(&st))
	}

	buf, err := json.Marshal(stamps)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(buf), nil
}

func (_stampChallengeConv) fromDB(v json.RawMessage) (*domain.StampChallenges, error) {

	conv := func(sc *_stampChallenge) (*domain.StampChallenge, error) {
		slotID, err := value.NewSlotID(sc.SlotID)
		if err != nil {
			return nil, err
		}
		cond, err := value.NewStampCondtion(value.StampConditionKind(sc.Condition))
		if err != nil {
			return nil, err
		}

		return &domain.StampChallenge{
			SlotID:    slotID,
			Condition: cond,
			UpdatedAt: sc.UpdatedAt,
		}, nil
	}

	var stamps []_stampChallenge
	if err := json.Unmarshal(v, &stamps); err != nil {
		return nil, err
	}

	var items []domain.StampChallenge
	for _, st := range stamps {
		dst, err := conv(&st)
		if err != nil {
			return nil, err
		}
		items = append(items, *dst)
	}

	return &domain.StampChallenges{Items: items}, nil
}

var watchEventConv _watchEventConv

type _watchEventConv struct{}

func (_watchEventConv) fromDB(v []WatchEvent) (*domain.WatchEvents, error) {

	conv := func(v *WatchEvent) (*domain.WatchEvent, error) {
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
		viewingSeconds, err := value.NewViewingPeriod(v.ViewingSeconds)
		if err != nil {
			return nil, err
		}
		return &domain.WatchEvent{
			TrackID:        trackID,
			TalkID:         talkID,
			SlotID:         slotID,
			ViewingSeconds: viewingSeconds,
			CreatedAt:      v.CreatedAt,
		}, nil
	}

	var items []domain.WatchEvent

	for _, ev := range v {
		dev, err := conv(&ev)
		if err != nil {
			return nil, err
		}
		items = append(items, *dev)
	}

	return &domain.WatchEvents{Items: items}, nil
}

func (_watchEventConv) toDB(confName value.ConfName, profileID value.ProfileID, v *domain.WatchEvents) ([]WatchEvent, error) {

	conv := func(dev *domain.WatchEvent) *WatchEvent {
		return &WatchEvent{
			ConferenceName: string(confName.Value()),
			ProfileID:      profileID.Value(),
			TrackID:        dev.TrackID.Value(),
			TalkID:         dev.TalkID.Value(),
			SlotID:         dev.SlotID.Value(),
			ViewingSeconds: dev.ViewingSeconds.Value(),
			CreatedAt:      dev.CreatedAt,
		}
	}

	var events []WatchEvent
	for _, ev := range v.Items {
		events = append(events, *conv(&ev))
	}

	return events, nil
}
