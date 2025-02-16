package dkui_test

import (
	"context"
	"net/url"
	"testing"

	"dreamkast-weaver/internal/dkui"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"

	"github.com/ServiceWeaver/weaver/weavertest"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	weaverConfig = `
	["dreamkast-weaver/internal/dkui/Service"]
	db_user = "user"
	db_password = "password"
	db_endpoint = "127.0.0.1"
	db_port = "13306"
	db_name = "test_dkui"
	`
	dbUrl = "mysql://user:password@127.0.0.1:13306/test_dkui"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	u, _ := url.Parse(dbUrl)
	db := dbmate.New(u)
	db.AutoDumpSchema = false

	mustNil(db.Drop())
	mustNil(db.CreateAndMigrate())
}

func teardown() {}

func TestDkUiServiceImpl_CreateViewEvent(t *testing.T) {
	domain.ChangeGuardSecondsForTest(0)
	domain.ChangeStampReadySecondsForTest(value.INTERVAL_SECONDS * 2)

	runner := weavertest.Local
	runner.Config = weaverConfig
	runner.Test(t, func(t *testing.T, svc dkui.Service) {
		ctx := context.Background()

		profile := dkui.Profile{
			ConfName: newConfName("cndt2023"),
			ID:       newProfileID(1),
		}
		req := dkui.CreateViewEventRequest{
			TrackID: newTrackID(2),
			TalkID:  newTalkID(3),
			SlotID:  newSlotID(1000),
		}

		// first time
		err := svc.CreateViewEvent(ctx, profile, req)
		assert.NoError(t, err)

		events, err := svc.ViewingEvents(ctx, profile)
		assert.NoError(t, err)

		stamps, err := svc.StampChallenges(ctx, profile)
		assert.NoError(t, err)

		assertViewEvents(t, events, 1000, value.INTERVAL_SECONDS)
		assert.Len(t, stamps.Items, 0)

		// second time
		err = svc.CreateViewEvent(ctx, profile, req)
		assert.NoError(t, err)

		events, err = svc.ViewingEvents(ctx, profile)
		assert.NoError(t, err)

		stamps, err = svc.StampChallenges(ctx, profile)
		assert.NoError(t, err)

		assertViewEvents(t, events, 1000, value.INTERVAL_SECONDS*2)
		assert.Len(t, stamps.Items, 1)
		assertStampCondition(t, stamps, 1000, "ready")

		// stamp
		err = svc.StampOnline(ctx, profile, req.SlotID)
		assert.NoError(t, err)

		events, err = svc.ViewingEvents(ctx, profile)
		assert.NoError(t, err)

		stamps, err = svc.StampChallenges(ctx, profile)
		assert.NoError(t, err)

		assertViewEvents(t, events, 1000, value.INTERVAL_SECONDS*2)
		assert.Len(t, stamps.Items, 1)
		assertStampCondition(t, stamps, 1000, "stamped")
	})
}

func TestDkUiServiceImpl_StampOnSite(t *testing.T) {
	runner := weavertest.Local
	runner.Config = weaverConfig
	runner.Test(t, func(t *testing.T, svc dkui.Service) {
		ctx := context.Background()

		profile := dkui.Profile{
			ConfName: newConfName("cndt2023"),
			ID:       newProfileID(1),
		}
		req := dkui.StampRequest{
			TrackID: newTrackID(2),
			TalkID:  newTalkID(3),
			SlotID:  newSlotID(1001),
		}

		err := svc.StampOnSite(ctx, profile, req)
		assert.NoError(t, err)

		slots, err := svc.ViewingEvents(ctx, profile)
		assert.NoError(t, err)

		stamps, err := svc.StampChallenges(ctx, profile)
		assert.NoError(t, err)

		assertViewEvents(t, slots, 1001, value.TALK_SECONDS)
		assertStampCondition(t, stamps, 1001, "stamped")
	})
}

func TestDkUiServiceImpl_ListTrackViewer(t *testing.T) {
	runner := weavertest.Local
	runner.Config = weaverConfig
	runner.Test(t, func(t *testing.T, svc dkui.Service) {
		ctx := context.Background()

		tna := newTrackName("A")
		tnb := newTrackName("B")
		tnc := newTrackName("C")
		tIDa := newTalkID(1)
		tIDb := newTalkID(2)
		tIDc := newTalkID(3)

		assert.NoError(t, svc.ViewTrack(ctx, newProfileID(731), tna, tIDa))
		assert.NoError(t, svc.ViewTrack(ctx, newProfileID(731), tna, tIDa))
		assert.NoError(t, svc.ViewTrack(ctx, newProfileID(732), tnb, tIDb))
		assert.NoError(t, svc.ViewTrack(ctx, newProfileID(733), tnc, tIDb))
		assert.NoError(t, svc.ViewTrack(ctx, newProfileID(734), tnc, tIDc))

		ans := map[value.TrackName]int{}
		ans[tna] = 1
		ans[tnb] = 1
		ans[tnc] = 2

		dvc, err := svc.ListViewerCounts(ctx, false)
		assert.NoError(t, err)

		for _, v := range dvc.Items {
			assert.Equal(t, ans[v.TrackName], v.Count)
		}
	})
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStampCondition(t *testing.T, stamps *domain.StampChallenges, slotID int32, status value.StampConditionKind) {
	t.Helper()
	var found bool
	for _, sc := range stamps.Items {
		if sc.SlotID.Value() == slotID {
			found = true
			assert.Equal(t, status, sc.Condition.Value())
		}
	}
	assert.True(t, found)
}

func assertViewEvents(t *testing.T, events *domain.ViewEvents, slotID int32, vt int32) {
	t.Helper()
	var total int32
	for _, ev := range events.Items {
		if ev.SlotID.Value() == slotID {
			total += ev.ViewingSeconds.Value()
		}
	}
	assert.Equal(t, vt, total)
}

func newProfile(confName value.ConferenceKind, profile int32) dkui.Profile {
	return dkui.Profile{
		ConfName: newConfName(confName),
		ID:       newProfileID(profile),
	}
}

func newConfName(v value.ConferenceKind) value.ConfName {
	o, err := value.NewConfName(v)
	mustNil(err)
	return o
}

func newProfileID(v int32) value.ProfileID {
	o, err := value.NewProfileID(v)
	mustNil(err)
	return o
}

func newTrackID(v int32) value.TrackID {
	o, err := value.NewTrackID(v)
	mustNil(err)
	return o
}

func newTalkID(v int32) value.TalkID {
	o, err := value.NewTalkID(v)
	mustNil(err)
	return o
}

func newSlotID(v int32) value.SlotID {
	o, err := value.NewSlotID(v)
	mustNil(err)
	return o
}

func newViewingSeconds(v int32) value.ViewingSeconds {
	o, err := value.NewViewingSeconds(v)
	mustNil(err)
	return o
}

func newStampCondition(v value.StampConditionKind) value.StampCondition {
	o, err := value.NewStampCondition(v)
	mustNil(err)
	return o
}

func newTrackName(v string) value.TrackName {
	o, err := value.NewTrackName(v)
	mustNil(err)
	return o
}
