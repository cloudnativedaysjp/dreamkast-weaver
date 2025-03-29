package application

import (
	"context"
	"net/url"
	"testing"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/stretchr/testify/assert"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/pkg/sqlhelper"
)

const (
	dkuiDBUrl = "mysql://user:password@127.0.0.1:13306/test_dkui"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	uu, _ := url.Parse(dkuiDBUrl)
	db := dbmate.New(uu)
	db.AutoDumpSchema = false
	db.MigrationsDir = []string{"../infrastructure/db/migrations"}

	mustNil(db.Drop())
	mustNil(db.CreateAndMigrate())
}

func teardown() {}

func TestStampRallyAppImpl_CreateViewEvent(t *testing.T) {
	dmodel.ChangeGuardSecondsForTest(0)
	dmodel.ChangeStampReadySecondsForTest(value.INTERVAL_SECONDS * 2)

	ctx := context.Background()

	sq, _ := sqlhelper.NewSqlHelper(&sqlhelper.SqlOption{
		User:     "user",
		Password: "password",
		Endpoint: "127.0.0.1",
		Port:     "13306",
		DbName:   "test_dkui",
	})

	svc := NewStampRallyApp(sq)

	profile := Profile{
		ConfName: newConfName("cndt2023"),
		ID:       newProfileID(1),
	}
	req := CreateViewEventRequest{
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
}

func TestStampRallyAppImpl_StampOnSite(t *testing.T) {

	ctx := context.Background()
	sq, _ := sqlhelper.NewSqlHelper(&sqlhelper.SqlOption{
		User:     "user",
		Password: "password",
		Endpoint: "127.0.0.1",
		Port:     "13306",
		DbName:   "test_dkui",
	})

	svc := NewStampRallyApp(sq)

	profile := Profile{
		ConfName: newConfName("cndt2023"),
		ID:       newProfileID(1),
	}
	req := StampRequest{
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
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStampCondition(t *testing.T, stamps *dmodel.StampChallenges, slotID int32, status value.StampConditionKind) {
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

func assertViewEvents(t *testing.T, events *dmodel.ViewEvents, slotID int32, vt int32) {
	t.Helper()
	var total int32
	for _, ev := range events.Items {
		if ev.SlotID.Value() == slotID {
			total += ev.ViewingSeconds.Value()
		}
	}
	assert.Equal(t, vt, total)
}

func newProfile(confName value.ConferenceKind, profile int32) Profile {
	return Profile{
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
