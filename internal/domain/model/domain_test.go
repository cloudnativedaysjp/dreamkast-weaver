package model_test

import (
	"net"
	"testing"
	"time"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"

	"github.com/stretchr/testify/assert"
)

var (
	svcDkui = dmodel.DkUiDomain{}
)

func TestDkUiService_CreateOnlineViewEvent(t *testing.T) {
	slotID := newSlotID(42)
	trackID := newTrackID(1)
	talkID := newTalkID(2)

	tests := []struct {
		name                      string
		given                     func() *dmodel.ViewEvents
		shouldStampChallengeAdded bool
	}{
		{
			name: "stamp condition fulfilled",
			given: func() *dmodel.ViewEvents {
				events := &dmodel.ViewEvents{}
				for i := 0; i < 9; i++ {
					ev := *dmodel.NewOnlineViewEvent(newTrackID(11), newTalkID(22), slotID)
					ev.CreatedAt = ev.CreatedAt.Add(-1 * (value.GUARD_SECONDS + 1) * time.Second)
					events = events.AddImmutable(ev)
				}
				return events
			},
			shouldStampChallengeAdded: true,
		},
		{
			name: "stamp condition not fulfilled",
			given: func() *dmodel.ViewEvents {
				events := &dmodel.ViewEvents{}
				for i := 0; i < 8; i++ {
					ev := *dmodel.NewOnlineViewEvent(newTrackID(11), newTalkID(22), slotID)
					ev.CreatedAt = ev.CreatedAt.Add(-1 * (value.GUARD_SECONDS + 1) * time.Second)
					events = events.AddImmutable(ev)
				}
				return events
			},
			shouldStampChallengeAdded: false,
		},
		{
			name: "first event",
			given: func() *dmodel.ViewEvents {
				return &dmodel.ViewEvents{}
			},
			shouldStampChallengeAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run("ok:"+tt.name, func(t *testing.T) {
			stamps := &dmodel.StampChallenges{}
			events := tt.given()
			evLen := len(events.Items)

			got, err := svcDkui.CreateOnlineViewEvent(trackID, talkID, slotID, stamps, events)

			assert.NoError(t, err)
			assert.Equal(t, trackID, got.TrackID)
			assert.Equal(t, talkID, got.TalkID)
			assert.Equal(t, slotID, got.SlotID)
			assert.Equal(t, value.ViewingSeconds120, got.ViewingSeconds)
			assert.Equal(t, evLen, len(events.Items))
			if tt.shouldStampChallengeAdded {
				assert.Equal(t, 1, len(stamps.Items))
				stamp := stamps.Items[0]
				assert.Equal(t, value.StampReady, stamp.Condition)
			} else {
				assert.Equal(t, 0, len(stamps.Items))
			}
		})
	}

	errTests := []struct {
		name  string
		given func() (*dmodel.ViewEvents, *dmodel.StampChallenges)
	}{
		{
			name: "too short request",
			given: func() (*dmodel.ViewEvents, *dmodel.StampChallenges) {
				events := &dmodel.ViewEvents{}
				ev := *dmodel.NewOnlineViewEvent(newTrackID(11), newTalkID(22), slotID)
				ev.CreatedAt = ev.CreatedAt.Add(-1 * (value.GUARD_SECONDS - 9) * time.Second)
				events = events.AddImmutable(ev)
				return events, &dmodel.StampChallenges{}
			},
		},
		{
			name: "nil given",
			given: func() (*dmodel.ViewEvents, *dmodel.StampChallenges) {
				return nil, nil
			},
		},
	}

	for _, tt := range errTests {
		t.Run("err:"+tt.name, func(t *testing.T) {
			events, stamps := tt.given()
			_, err := svcDkui.CreateOnlineViewEvent(trackID, talkID, slotID, stamps, events)
			assert.Error(t, err)
		})
	}
}

func TestDkUiService_StampOnline(t *testing.T) {
	slotID := newSlotID(42)

	t.Run("ok", func(t *testing.T) {
		stamps := &dmodel.StampChallenges{Items: []dmodel.StampChallenge{
			*dmodel.NewStampChallenge(newSlotID(41)),
			*dmodel.NewStampChallenge(newSlotID(42)),
			*dmodel.NewStampChallenge(newSlotID(43)),
		}}

		err := svcDkui.StampOnline(slotID, stamps)

		assert.NoError(t, err)
		for _, stamp := range stamps.Items {
			if stamp.SlotID == slotID {
				assert.Equal(t, value.StampStamped, stamp.Condition)
			} else {
				assert.Equal(t, value.StampSkipped, stamp.Condition)
			}
		}
	})

	errTests := []struct {
		name  string
		given func() *dmodel.StampChallenges
	}{
		{
			name: "ready stamp not found",
			given: func() *dmodel.StampChallenges {
				return &dmodel.StampChallenges{Items: []dmodel.StampChallenge{
					*dmodel.NewStampChallenge(newSlotID(41)),
					*dmodel.NewStampChallenge(newSlotID(43)),
				}}
			},
		},
		{
			name: "nil given",
			given: func() *dmodel.StampChallenges {
				return nil
			},
		},
	}

	for _, tt := range errTests {
		t.Run("err:"+tt.name, func(t *testing.T) {
			stamps := tt.given()

			err := svcDkui.StampOnline(slotID, stamps)

			assert.Error(t, err)
			if stamps != nil {
				for _, stamp := range stamps.Items {
					assert.Equal(t, value.StampReady, stamp.Condition)
				}
			}
		})
	}
}

func TestDkUiService_StampOnSite(t *testing.T) {
	slotID := newSlotID(42)
	trackID := newTrackID(1)
	talkID := newTalkID(2)

	tests := []struct {
		name  string
		given func() *dmodel.StampChallenges
	}{
		{
			name: "stamp not exist",
			given: func() *dmodel.StampChallenges {
				return &dmodel.StampChallenges{Items: []dmodel.StampChallenge{
					*dmodel.NewStampChallenge(newSlotID(41)),
					*dmodel.NewStampChallenge(newSlotID(43)),
				}}
			},
		},
		{
			name: "ready stamp exists",
			given: func() *dmodel.StampChallenges {
				return &dmodel.StampChallenges{Items: []dmodel.StampChallenge{
					*dmodel.NewStampChallenge(newSlotID(41)),
					*dmodel.NewStampChallenge(newSlotID(42)),
					*dmodel.NewStampChallenge(newSlotID(43)),
				}}
			},
		},
	}

	for _, tt := range tests {
		t.Run("ok"+tt.name, func(t *testing.T) {
			stamps := tt.given()

			got, err := svcDkui.StampOnSite(trackID, talkID, slotID, stamps)

			assert.NoError(t, err)
			assert.Equal(t, trackID, got.TrackID)
			assert.Equal(t, talkID, got.TalkID)
			assert.Equal(t, slotID, got.SlotID)
			assert.Equal(t, value.ViewingSeconds2400, got.ViewingSeconds)
			for _, stamp := range stamps.Items {
				if stamp.SlotID == slotID {
					assert.Equal(t, value.StampStamped, stamp.Condition)
				} else {
					assert.Equal(t, value.StampReady, stamp.Condition)
				}
			}
		})
	}

	errTests := []struct {
		name  string
		given func() *dmodel.StampChallenges
	}{
		{
			name: "already stamped",
			given: func() *dmodel.StampChallenges {
				sc := dmodel.NewStampChallenge(newSlotID(42))
				sc.Stamp()
				return &dmodel.StampChallenges{Items: []dmodel.StampChallenge{
					*dmodel.NewStampChallenge(newSlotID(41)),
					*sc,
					*dmodel.NewStampChallenge(newSlotID(43)),
				}}
			},
		},
		{
			name: "nil given",
			given: func() *dmodel.StampChallenges {
				return nil
			},
		},
	}

	for _, tt := range errTests {
		t.Run("err"+tt.name, func(t *testing.T) {
			stamps := tt.given()

			_, err := svcDkui.StampOnSite(trackID, talkID, slotID, stamps)

			assert.Error(t, err)
			if stamps != nil {
				for _, stamp := range stamps.Items {
					if stamp.SlotID == slotID {
						assert.Equal(t, value.StampStamped, stamp.Condition)
					} else {
						assert.Equal(t, value.StampReady, stamp.Condition)
					}
				}
			}
		})
	}
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func newSlotID(v int32) value.SlotID {
	id, err := value.NewSlotID(v)
	mustNil(err)
	return id
}

func newTrackID(v int32) value.TrackID {
	id, err := value.NewTrackID(v)
	mustNil(err)
	return id
}

func newTalkID(v int32) value.TalkID {
	id, err := value.NewTalkID(v)
	mustNil(err)
	return id
}

func TestCfpDomain_TallyCfpVotes(t *testing.T) {
	ss := newSpanSeconds(value.SPAN_SECONDS)
	tn := time.Unix(time.Now().Unix()/value.SPAN_SECONDS*value.SPAN_SECONDS, 0)
	id := newTalkID(1)
	ip := net.ParseIP("192.0.2.1")

	tests := []struct {
		name   string
		given  func() (cvs *dmodel.CfpVotes)
		counts map[value.TalkID]int32
	}{
		{
			name: "votes within a span are summarized",
			given: func() (cvs *dmodel.CfpVotes) {
				cvs = &dmodel.CfpVotes{Items: []dmodel.CfpVote{
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn,
					},
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn.Add((value.SPAN_SECONDS - 1) * time.Second),
					},
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn.Add(value.SPAN_SECONDS * time.Second),
					},
				}}
				return cvs
			},
			counts: map[value.TalkID]int32{
				id: 2,
			},
		},
		{
			name: "different IPs or IDs are different votes",
			given: func() (cvs *dmodel.CfpVotes) {
				cvs = &dmodel.CfpVotes{Items: []dmodel.CfpVote{
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn,
					},
					{
						TalkID:    newTalkID(2),
						ClientIp:  ip,
						CreatedAt: tn.Add((value.SPAN_SECONDS - 3) * time.Second),
					},
					{
						TalkID:    id,
						ClientIp:  net.ParseIP("192.0.2.2"),
						CreatedAt: tn.Add((value.SPAN_SECONDS - 1) * time.Second),
					},
				}}
				return cvs
			},
			counts: map[value.TalkID]int32{
				id:           2,
				newTalkID(2): 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run("ok:"+tt.name, func(t *testing.T) {
			cvs := tt.given()

			got := cvs.Tally(ss)

			for _, v := range got {
				assert.Equal(t, tt.counts[v.TalkID], v.Count)
			}
		})
	}
}

func newSpanSeconds(v int32) value.SpanSeconds {
	ss, err := value.NewSpanSeconds(&v)
	mustNil(err)
	return ss
}
