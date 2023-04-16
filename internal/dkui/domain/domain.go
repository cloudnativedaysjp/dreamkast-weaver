package domain

import (
	"fmt"
	"time"

	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/stacktrace"
)

var (
	viewEventGuardSeconds = value.GUARD_SECONDS
	stampReadySeconds     = value.STAMP_READY_SECONDS
	jst                   *time.Location
)

var (
	ErrTooShortRequest = derrors.NewUserError("too short requests")
	ErrStampNotReady   = derrors.NewUserError("stamp is not ready")
	ErrAlreadyStamped  = derrors.NewUserError("already stamped")
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

func ChangeGuardSecondsForTest(sec int) {
	viewEventGuardSeconds = sec
}

func ChangeStampReadySecondsForTest(sec int) {
	stampReadySeconds = sec
}

func nowJST() time.Time {
	return time.Now().In(jst)
}

type DkUiDomain struct{}

func (DkUiDomain) CreateOnlineViewEvent(
	trackID value.TrackID,
	talkID value.TalkID,
	slotID value.SlotID,
	stamps *StampChallenges,
	events *ViewEvents) (*ViewEvent, error) {
	if stamps == nil || events == nil {
		return nil, stacktrace.With(fmt.Errorf("missing required params"))
	}
	ev := NewOnlineViewEvent(trackID, talkID, slotID)

	lastCreatedAt := events.LastCreated()
	if ev.CreatedAt.Sub(lastCreatedAt) < time.Duration(viewEventGuardSeconds)*time.Second {
		return nil, ErrTooShortRequest
	}

	stamps.MakeReadyIfFulfilled(slotID, events.AddImmutable(*ev))
	return ev, nil
}

func (DkUiDomain) StampOnline(
	slotID value.SlotID,
	stamps *StampChallenges) error {
	if stamps == nil {
		return stacktrace.With(fmt.Errorf("missing required params"))
	}

	return stamps.StampIfReady(slotID)
}

func (DkUiDomain) StampOnSite(
	trackID value.TrackID,
	talkID value.TalkID,
	slotID value.SlotID,
	stamps *StampChallenges) (*ViewEvent, error) {
	if stamps == nil {
		return nil, stacktrace.With(fmt.Errorf("missing required params"))
	}

	if err := stamps.ForceStamp(slotID); err != nil {
		return nil, err
	}
	return NewOnSiteViewEvent(trackID, talkID, slotID), nil
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

func (scs *StampChallenges) MakeReadyIfFulfilled(slotID value.SlotID, evs *ViewEvents) {
	if evs.IsFulfilled(slotID) {
		scs.setReadyChallenge(slotID)
	}
}

func (scs *StampChallenges) StampIfReady(slotID value.SlotID) error {
	sc := scs.Get(slotID)
	if sc == nil || sc.Condition != value.StampReady {
		return ErrStampNotReady
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
	sc := scs.Get(slotID)
	if sc == nil {
		scs.setReadyChallenge(slotID)
	}
	if sc != nil && sc.Condition == value.StampStamped {
		return ErrAlreadyStamped
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
	sc := scs.Get(slotID)
	if sc == nil {
		scs.Items = append(scs.Items, *NewStampChallenge(slotID))
	}
}

func (scs *StampChallenges) Get(slotID value.SlotID) *StampChallenge {
	for _, p := range scs.Items {
		sc := p
		if sc.SlotID == slotID {
			return &sc
		}
	}
	return nil
}

type ViewEvent struct {
	TrackID        value.TrackID
	TalkID         value.TalkID
	SlotID         value.SlotID
	ViewingSeconds value.ViewingSeconds
	CreatedAt      time.Time
}

func NewOnlineViewEvent(trackID value.TrackID, talkID value.TalkID, slotID value.SlotID) *ViewEvent {
	return &ViewEvent{
		TrackID:        trackID,
		TalkID:         talkID,
		SlotID:         slotID,
		ViewingSeconds: value.ViewingSeconds120,
		CreatedAt:      nowJST(),
	}
}

func NewOnSiteViewEvent(trackID value.TrackID, talkID value.TalkID, slotID value.SlotID) *ViewEvent {
	return &ViewEvent{
		TrackID:        trackID,
		TalkID:         talkID,
		SlotID:         slotID,
		ViewingSeconds: value.ViewingSeconds2400,
		CreatedAt:      nowJST(),
	}
}

type ViewEvents struct {
	Items []ViewEvent
}

func (evs *ViewEvents) IsFulfilled(slotID value.SlotID) bool {
	var total int32
	for _, ev := range evs.Items {
		if ev.SlotID == slotID {
			total += ev.ViewingSeconds.Value()
		}
	}
	return total >= int32(stampReadySeconds)
}

func (evs *ViewEvents) LastCreated() time.Time {
	var lastTime time.Time
	for _, ev := range evs.Items {
		if ev.CreatedAt.After(lastTime) {
			lastTime = ev.CreatedAt
		}
	}
	return lastTime
}

func (evs *ViewEvents) ViewingSeconds() map[value.SlotID]int32 {
	res := map[value.SlotID]int32{}

	for _, ev := range evs.Items {
		res[ev.SlotID] += ev.ViewingSeconds.Value()
	}
	return res
}

func (evs *ViewEvents) AddImmutable(ev ViewEvent) *ViewEvents {
	events := make([]ViewEvent, len(evs.Items)+1)
	events[0] = ev
	copy(events[1:], evs.Items)
	return &ViewEvents{
		Items: events,
	}
}
