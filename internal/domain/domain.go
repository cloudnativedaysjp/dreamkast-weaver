package domain

import (
	"errors"
	"net"
	"time"

	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/stacktrace"

	"github.com/ServiceWeaver/weaver"
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

	ErrMissingParams = errors.New("missing required params")
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

func ChangeStampReadySecondsForTest(sec int32) {
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
		return nil, stacktrace.With(ErrMissingParams)
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
		return stacktrace.With(ErrMissingParams)
	}

	return stamps.StampIfReady(slotID)
}

func (DkUiDomain) StampOnSite(
	trackID value.TrackID,
	talkID value.TalkID,
	slotID value.SlotID,
	stamps *StampChallenges) (*ViewEvent, error) {
	if stamps == nil {
		return nil, stacktrace.With(ErrMissingParams)
	}

	if err := stamps.ForceStamp(slotID); err != nil {
		return nil, err
	}
	return NewOnSiteViewEvent(trackID, talkID, slotID), nil
}

type StampChallenge struct {
	weaver.AutoMarshal
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
	weaver.AutoMarshal
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
	weaver.AutoMarshal
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
	weaver.AutoMarshal
	Items []ViewEvent
}

func (evs *ViewEvents) IsFulfilled(slotID value.SlotID) bool {
	var total int32
	for _, ev := range evs.Items {
		if ev.SlotID == slotID {
			total += ev.ViewingSeconds.Value()
		}
	}
	return total >= stampReadySeconds
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

func NewViewerCount(tn value.TrackName, count int) *ViewerCount {
	return &ViewerCount{
		TrackName: tn,
		Count:     count,
	}
}

type ViewerCount struct {
	weaver.AutoMarshal
	TrackName value.TrackName
	Count     int
}

type ViewerCounts struct {
	weaver.AutoMarshal
	Items []ViewerCount
}

type TrackViewer struct {
	weaver.AutoMarshal
	CreatedAt time.Time
	TrackName value.TrackName
	ProfileID value.ProfileID
}

type TrackViewers struct {
	weaver.AutoMarshal
	Items []TrackViewer
}

func (v *TrackViewers) GetViewerCounts() ViewerCounts {
	aa := map[value.TrackName]int{}
	for _, n := range value.TrackNames() {
		aa[n] = 0
	}

	type key struct {
		tn        value.TrackName
		profileID value.ProfileID
	}
	counted := map[key]bool{}

	for _, v := range v.Items {
		k := key{
			tn:        v.TrackName,
			profileID: v.ProfileID,
		}
		if _, isThere := counted[k]; isThere {
			continue
		}
		counted[k] = true
		aa[v.TrackName]++
	}

	var items []ViewerCount
	for k, v := range aa {
		items = append(items, ViewerCount{
			TrackName: k,
			Count:     v,
		})
	}

	return ViewerCounts{
		Items: items,
	}
}

type CfpVote struct {
	weaver.AutoMarshal
	TalkID    value.TalkID
	ClientIp  net.IP
	CreatedAt time.Time
}

type CfpVotes struct {
	weaver.AutoMarshal
	Items []CfpVote
}

type VoteCount struct {
	weaver.AutoMarshal
	TalkID value.TalkID
	Count  int
}

type CfpDomain struct{}

func (cd *CfpDomain) TallyCfpVotes(cfpVotes *CfpVotes, spanSeconds value.SpanSeconds) []*VoteCount {
	type key struct {
		talkId    int32
		ip        string
		timeFrame int64
	}
	voted := map[key]bool{}
	counts := map[value.TalkID]int{}
	for _, v := range cfpVotes.Items {
		k := key{
			talkId:    v.TalkID.Value(),
			ip:        v.ClientIp.String(),
			timeFrame: v.CreatedAt.Unix() / int64(spanSeconds.Value()),
		}
		if _, isThere := voted[k]; isThere {
			continue
		}
		voted[k] = true
		counts[v.TalkID]++
	}

	var resp []*VoteCount
	for talkID, count := range counts {
		resp = append(resp, &VoteCount{
			TalkID: talkID,
			Count:  count,
		})
	}
	return resp
}
