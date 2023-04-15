package dkui_test

import (
	"context"
	"dreamkast-weaver/internal/dkui"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/graph/model"
	"net/url"
	"testing"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/weavertest"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/stretchr/testify/assert"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
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

	mustNil(db.Drop())
	mustNil(db.CreateAndMigrate())
}

func teardown() {}

func TestDkUiServiceImpl_CreateWatchEvent(t *testing.T) {
	domain.ChangeGuardSecondsForTest(0)
	domain.ChangeStampReadySecondsForTest(value.INTERVAL_SECONDS * 2)

	ctx := context.Background()
	root := weavertest.Init(ctx, t, weavertest.Options{
		SingleProcess: true,
		Config:        weaverConfig,
	})
	svc, err := weaver.Get[dkui.Service](root)
	mustNil(err)

	req := model.CreateWatchEventInput{
		ConfName:  "cndf2023",
		ProfileID: 1,
		TrackID:   2,
		TalkID:    3,
		SlotID:    1000,
	}

	// first time
	err = svc.CreateWatchEvent(ctx, req)
	assert.Nil(t, err)

	slots, err := svc.ViewingSlots(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	stamps, err := svc.StampChallenges(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	assertViewingTime(t, slots, 1000, value.INTERVAL_SECONDS)
	assert.Len(t, stamps, 0)

	// second time
	err = svc.CreateWatchEvent(ctx, req)
	assert.Nil(t, err)

	slots, err = svc.ViewingSlots(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	stamps, err = svc.StampChallenges(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	assertViewingTime(t, slots, 1000, value.INTERVAL_SECONDS*2)
	assert.Len(t, stamps, 1)
	assertStampCondition(t, stamps, 1000, "ready")

	// stamp
	stampReq := model.StampOnlineInput{
		ConfName:  "cndf2023",
		ProfileID: 1,
		SlotID:    1000,
	}
	err = svc.StampOnline(ctx, stampReq)
	assert.Nil(t, err)

	slots, err = svc.ViewingSlots(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	stamps, err = svc.StampChallenges(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	assertViewingTime(t, slots, 1000, value.INTERVAL_SECONDS*2)
	assert.Len(t, stamps, 1)
	assertStampCondition(t, stamps, 1000, "stamped")
}

func TestDkUiServiceImpl_StampOnSite(t *testing.T) {

	ctx := context.Background()
	root := weavertest.Init(ctx, t, weavertest.Options{
		SingleProcess: true,
		Config:        weaverConfig,
	})
	svc, err := weaver.Get[dkui.Service](root)
	mustNil(err)

	req := model.StampOnSiteInput{
		ConfName:  "cndf2023",
		ProfileID: 1,
		TrackID:   2,
		TalkID:    3,
		SlotID:    1001,
	}

	err = svc.StampOnSite(ctx, req)
	assert.Nil(t, err)

	slots, err := svc.ViewingSlots(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	stamps, err := svc.StampChallenges(ctx, "cndf2023", 1)
	assert.Nil(t, err)

	assertViewingTime(t, slots, 1001, value.TALK_SECONDS)
	assertStampCondition(t, stamps, 1001, "stamped")
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStampCondition(t *testing.T, stamps []*model.StampChallenge, slotID int, status string) {
	t.Helper()
	var found bool
	for _, sc := range stamps {
		if sc.SlotID == slotID {
			found = true
			assert.Equal(t, status, sc.Condition.String())
		}
	}
	assert.True(t, found)
}

func assertViewingTime(t *testing.T, slots []*model.ViewingSlot, slotID int, vt int) {
	t.Helper()
	var found bool
	for _, s := range slots {
		if s.SlotID == slotID {
			found = true
			assert.Equal(t, vt, s.ViewingTime)
		}
	}
	assert.True(t, found)
}
