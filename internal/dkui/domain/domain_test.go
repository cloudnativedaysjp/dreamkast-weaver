package domain_test

import (
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"testing"
	"time"
)

func newSlotID(v int32) value.SlotID {
	id, _ := value.NewSlotID(v)
	return id
}

func newTrackID(v int32) value.TrackID {
	id, _ := value.NewTrackID(v)
	return id
}

func newTalkID(v int32) value.TalkID {
	id, _ := value.NewTalkID(v)
	return id
}

var (
	svc = domain.DkUiService{}
)

func TestDkUiService_CreateOnlineWatchEvent(t *testing.T) {

	slotID := newSlotID(42)
	trackID := newTrackID(1)
	talkID := newTalkID(2)

	tests := []struct {
		name                      string
		given                     func() *domain.WatchEvents
		shouldStampChallengeAdded bool
	}{
		{
			name: "stamp condition fulfilled",
			given: func() *domain.WatchEvents {
				events := &domain.WatchEvents{}
				for i := 0; i < 9; i++ {
					ev := *domain.NewOnlineWatchEvent(newTrackID(11), newTalkID(22), slotID)
					ev.CreatedAt = ev.CreatedAt.Add(time.Duration(-1 * (value.GUARD_SECONDS + 1) * time.Second))
					events = events.AddImmutable(ev)
				}
				return events
			},
			shouldStampChallengeAdded: true,
		},
		{
			name: "stamp condition not fulfilled",
			given: func() *domain.WatchEvents {
				events := &domain.WatchEvents{}
				for i := 0; i < 8; i++ {
					ev := *domain.NewOnlineWatchEvent(newTrackID(11), newTalkID(22), slotID)
					ev.CreatedAt = ev.CreatedAt.Add(time.Duration(-1 * (value.GUARD_SECONDS + 1) * time.Second))
					events = events.AddImmutable(ev)
				}
				return events
			},
			shouldStampChallengeAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run("ok:"+tt.name, func(t *testing.T) {
			stamps := &domain.StampChallenges{}
			events := tt.given()
			evLen := len(events.Items)

			got, err := svc.CreateOnlineWatchEvent(trackID, talkID, slotID, stamps, events)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.TrackID != trackID {
				t.Errorf("not equal: want=%#v, got=%#v", trackID, got.TrackID)
			}
			if got.TalkID != talkID {
				t.Errorf("not equal: want=%#v, got=%#v", talkID, got.TalkID)
			}
			if got.SlotID != slotID {
				t.Errorf("not equal: want=%#v, got=%#v", slotID, got.SlotID)
			}
			if got.ViewingSeconds != value.ViewingSeconds120 {
				t.Errorf("not equal: want=%#v, got=%#v", value.ViewingSeconds120, got.ViewingSeconds)
			}
			if len(events.Items) != evLen {
				t.Errorf("events mutated: got=%#v", events)
			}
			if tt.shouldStampChallengeAdded {
				if len(stamps.Items) == 0 {
					t.Fatalf("stamp is not added")
				}
				stamp := stamps.Items[0]
				if stamp.Condition != value.StampReady {
					t.Fatalf("added stamp is not in ready condition")
				}
			} else {
				if len(stamps.Items) != 0 {
					t.Fatalf("stamp added unexpectedly")
				}
			}
		})
	}

	errTests := []struct {
		name  string
		given func() *domain.WatchEvents
	}{
		{
			name: "too short request",
			given: func() *domain.WatchEvents {
				events := &domain.WatchEvents{}
				ev := *domain.NewOnlineWatchEvent(newTrackID(11), newTalkID(22), slotID)
				ev.CreatedAt = ev.CreatedAt.Add(time.Duration(-1 * (value.GUARD_SECONDS - 9) * time.Second))
				events = events.AddImmutable(ev)
				return events
			},
		},
	}

	for _, tt := range errTests {
		t.Run("err:"+tt.name, func(t *testing.T) {
			stamps := &domain.StampChallenges{}
			events := tt.given()

			_, err := svc.CreateOnlineWatchEvent(trackID, talkID, slotID, stamps, events)
			if err == nil {
				t.Errorf("error not raised")
			}
		})
	}

}

func TestDkUiService_StampOnline(t *testing.T) {

	slotID := newSlotID(42)

	t.Run("ok", func(t *testing.T) {
		stamps := &domain.StampChallenges{[]domain.StampChallenge{
			*domain.NewStampChallenge(newSlotID(41)),
			*domain.NewStampChallenge(newSlotID(42)),
			*domain.NewStampChallenge(newSlotID(43)),
		}}

		if err := svc.StampOnline(slotID, stamps); err != nil {
			t.Fatal("unexpected error")
		}

		for _, stamp := range stamps.Items {
			if stamp.SlotID == slotID {
				if stamp.Condition != value.StampStamped {
					t.Errorf("not equal: want=%#v, got=%#v", value.StampStamped, stamp.Condition)
				}
			} else {
				if stamp.Condition != value.StampSkipped {
					t.Errorf("not equal: want=%#v, got=%#v", value.StampSkipped, stamp.Condition)
				}
			}
		}
	})

	errTests := []struct {
		name  string
		given func() *domain.StampChallenges
	}{
		{
			name: "ready stamp not found",
			given: func() *domain.StampChallenges {
				return &domain.StampChallenges{[]domain.StampChallenge{
					*domain.NewStampChallenge(newSlotID(41)),
					*domain.NewStampChallenge(newSlotID(43)),
				}}
			},
		},
	}

	for _, tt := range errTests {
		t.Run("err:"+tt.name, func(t *testing.T) {
			stamps := tt.given()
			err := svc.StampOnline(slotID, stamps)
			if err == nil {
				t.Error("error not raised")
			}
		})
	}
}
