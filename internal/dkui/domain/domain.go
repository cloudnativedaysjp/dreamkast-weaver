package domain

import (
	"dreamkast-weaver/internal/dkui/value"
	"fmt"
	"time"
)

var (
	viewEventGuardSeconds = value.GUARD_SECONDS
	jst                   *time.Location
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

func ChangeGuardSecondsForTest(guardSeconds int) {
	viewEventGuardSeconds = guardSeconds
}

func nowJST() time.Time {
	return time.Now().In(jst)
}

type DkUiDomain struct{}

func (DkUiDomain) CreateOnlineWatchEvent(
	trackID value.TrackID,
	talkID value.TalkID,
	slotID value.SlotID,
	stamps *StampChallenges,
	events *WatchEvents) (*WatchEvent, error) {
	if stamps == nil || events == nil {
		return nil, fmt.Errorf("missing required params")
	}

	ev := NewOnlineWatchEvent(trackID, talkID, slotID)

	lastCreatedAt := events.LastCreated()
	if ev.CreatedAt.Sub(lastCreatedAt) < time.Duration(viewEventGuardSeconds)*time.Second {
		return nil, fmt.Errorf("too short requests")
	}

	stamps.MakeReadyIfFulfilled(slotID, events.AddImmutable(*ev))
	return ev, nil
}

func (DkUiDomain) StampOnline(
	slotID value.SlotID,
	stamps *StampChallenges) error {
	if stamps == nil {
		return fmt.Errorf("missing required params")
	}

	return stamps.StampIfReady(slotID)
}

func (DkUiDomain) StampOnSite(
	trackID value.TrackID,
	talkID value.TalkID,
	slotID value.SlotID,
	stamps *StampChallenges) (*WatchEvent, error) {
	if stamps == nil {
		return nil, fmt.Errorf("missing required params")
	}

	if err := stamps.ForceStamp(slotID); err != nil {
		return nil, err
	}
	return NewOnSiteWatchEvent(trackID, talkID, slotID), nil
}

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
	sc.UpdatedAt = nowJST()
}

func (sc *StampChallenge) Skip() {
	sc.Condition = value.StampSkipped
	sc.UpdatedAt = nowJST()
}

type StampChallenges struct {
	Items []StampChallenge
}

func (scs *StampChallenges) MakeReadyIfFulfilled(slotID value.SlotID, evs *WatchEvents) {
	if evs.IsFulfilled(slotID) {
		scs.setReadyChallenge(slotID)
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

	for i, sc := range scs.Items {
		if sc.SlotID == slotID {
			sc.Stamp()
		}
		if sc.SlotID != slotID && sc.Condition == value.StampReady {
			sc.Skip()
		}
		scs.Items[i] = sc
	}
	return nil
}

func (scs *StampChallenges) ForceStamp(slotID value.SlotID) error {
	var tgt *StampChallenge
	for _, p := range scs.Items {
		sc := p
		if sc.SlotID == slotID {
			tgt = &sc
		}
	}
	if tgt == nil {
		scs.setReadyChallenge(slotID)
	}
	if tgt != nil && tgt.Condition == value.StampStamped {
		return fmt.Errorf("already stamped: slotID=%v", slotID)
	}

	for i, sc := range scs.Items {
		if sc.SlotID == slotID {
			sc.Stamp()
			scs.Items[i] = sc
		}
	}
	return nil
}

func (scs *StampChallenges) setReadyChallenge(slotID value.SlotID) {
	scs.Items = append(scs.Items, *NewStampChallenge(slotID))
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
		CreatedAt:      nowJST(),
	}
}

func NewOnSiteWatchEvent(trackID value.TrackID, talkID value.TalkID, slotID value.SlotID) *WatchEvent {
	return &WatchEvent{
		TrackID:        trackID,
		TalkID:         talkID,
		SlotID:         slotID,
		ViewingSeconds: value.ViewingSeconds2400,
		CreatedAt:      nowJST(),
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

func (evs *WatchEvents) ViewingSeconds() map[value.SlotID]int32 {
	res := map[value.SlotID]int32{}

	for _, ev := range evs.Items {
		res[ev.SlotID] += int32(ev.ViewingSeconds.Value())
	}
	return res
}

func (evs *WatchEvents) AddImmutable(ev WatchEvent) *WatchEvents {
	events := make([]WatchEvent, len(evs.Items)+1)
	events[0] = ev
	copy(events[1:], evs.Items)
	return &WatchEvents{
		Items: events,
	}
}
