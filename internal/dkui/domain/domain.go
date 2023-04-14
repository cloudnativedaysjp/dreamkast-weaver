package domain

import (
	"context"
	"dreamkast-weaver/internal/dkui/value"
	"fmt"
	"time"
)

type StampChallenge struct {
	SlotID    value.SlotID
	Condition value.StampCondition
	UpdatedAt time.Time
}

func NewStampChallenge(slotID value.SlotID) *StampChallenge {
	return &StampChallenge{
		SlotID:    slotID,
		Condition: value.StampReady,
		UpdatedAt: time.Time{},
	}
}

func (sc *StampChallenge) Stamp() {
	sc.Condition = value.StampStamped
	sc.UpdatedAt = time.Now()
}

func (sc *StampChallenge) Skip() {
	sc.Condition = value.StampSkipped
	sc.UpdatedAt = time.Now()
}

type StampChallenges struct {
	Items []StampChallenge
}

func (scs *StampChallenges) SetReadyIfFulfilled(slotID value.SlotID, evs *WatchEvents) {
	if evs.IsFulfilled(slotID) {
		fmt.Printf("fulfilled\n")
		scs.setReady(slotID)
	}
}

func (scs *StampChallenges) StampIfReady(slotID value.SlotID) error {
	var tgt *StampChallenge
	for _, sc := range scs.Items {
		if sc.SlotID == slotID {
			tgt = &sc
		}
	}
	if tgt == nil || tgt.Condition != value.StampReady {
		return fmt.Errorf("stamp is not ready: slotID=%v", slotID)
	}

	for _, sc := range scs.Items {
		if sc.SlotID == slotID {
			sc.Stamp()
		}
		if sc.SlotID != slotID && sc.Condition == value.StampReady {
			sc.Skip()
		}
	}
	return nil
}

func (scs *StampChallenges) ForceStamp(slotID value.SlotID) error {
	var tgt *StampChallenge
	for _, sc := range scs.Items {
		if sc.SlotID == slotID {
			tgt = &sc
		}
	}
	if tgt == nil {
		scs.setReady(slotID)
	}
	if tgt != nil && tgt.Condition == value.StampStamped {
		return fmt.Errorf("already stamped: slotID=%v", slotID)
	}

	for _, sc := range scs.Items {
		if sc.SlotID == slotID {
			sc.Stamp()
		}
		if sc.SlotID != slotID && sc.Condition == value.StampReady {
			sc.Skip()
		}
	}
	return nil
}

func (scs *StampChallenges) setReady(slotID value.SlotID) {
	scs.Items = append(scs.Items, *NewStampChallenge(slotID))
	fmt.Printf("scs.Items => %+v\n", scs.Items)
}

type WatchEvent struct {
	TrackID        value.TrackID
	TalkID         value.TalkID
	SlotID         value.SlotID
	ViewingSeconds value.ViewingSeconds
	CreatedAt      time.Time
}

func NewOnlineWatchEvent(trackID value.TrackID, talkID value.TalkID, slotID value.SlotID) *WatchEvent {
	return &WatchEvent{
		TrackID:        trackID,
		TalkID:         talkID,
		SlotID:         slotID,
		ViewingSeconds: value.ViewingSeconds120,
		CreatedAt:      time.Now(),
	}
}

func NewOnSiteWatchEvent(trackID value.TrackID, talkID value.TalkID, slotID value.SlotID) *WatchEvent {
	return &WatchEvent{
		TrackID:        trackID,
		TalkID:         talkID,
		SlotID:         slotID,
		ViewingSeconds: value.ViewingSeconds2400,
		CreatedAt:      time.Now(),
	}
}

type WatchEvents struct {
	Items []WatchEvent
}

func (evs *WatchEvents) IsFulfilled(slotID value.SlotID) bool {
	var total int32
	for _, ev := range evs.Items {
		if ev.SlotID == slotID {
			total += ev.ViewingSeconds.Value()
		}
	}
	return total >= value.STAMP_READY_SECONDS
}

func (evs *WatchEvents) LastCreated() time.Time {
	var lastTime time.Time
	for _, ev := range evs.Items {
		if ev.CreatedAt.After(lastTime) {
			lastTime = ev.CreatedAt
		}
	}
	return lastTime
}

func (evs *WatchEvents) AddImmutable(ev WatchEvent) *WatchEvents {
	events := make([]WatchEvent, len(evs.Items)+1)
	events[0] = ev
	copy(events[1:], evs.Items)
	return &WatchEvents{
		Items: events,
	}
}

type DkUiRepo interface {
	ListWatchEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (WatchEvents, error)
	InsertWatchEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev WatchEvent) error

	GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (StampChallenges, error)
	UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs StampChallenges) error
}

type DkUiService struct{}

func (DkUiService) CreateOnlineWatchEvent(
	trackID value.TrackID,
	talkID value.TalkID,
	slotID value.SlotID,
	stamps *StampChallenges,
	events *WatchEvents) (*WatchEvent, error) {

	ev := NewOnlineWatchEvent(trackID, talkID, slotID)
	lastCreatedAt := events.LastCreated()
	if ev.CreatedAt.Sub(lastCreatedAt) < value.GUARD_SECONDS*time.Second {
		return nil, fmt.Errorf("too short requests")
	}

	stamps.SetReadyIfFulfilled(slotID, events.AddImmutable(*ev))
	return ev, nil
}
