package dkui_test

import (
	"context"
	"dreamkast-weaver/internal/dkui"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/sqlhelper"
	"net/url"
	"testing"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/stretchr/testify/assert"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
)

const (
	testDB = "test_dkui"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	u, _ := url.Parse("mysql://user:password@127.0.0.1:13306/" + testDB)
	db := dbmate.New(u)

	mustNil(db.Drop())
	mustNil(db.CreateAndMigrate())
}

func teardown() {

}

func assertStampStatus(t *testing.T, st *dkui.StatusResponse, slotID int, status string) {
	t.Helper()
	var found bool
	for _, sc := range st.StampChallenges {
		if sc.SlotID == slotID {
			found = true
			assert.Equal(t, status, sc.Condition)
		}
	}
	assert.True(t, found)
}

func TestDkUiServiceImpl_CreateWatchEvent(t *testing.T) {
	domain.ChangeGuardSecondsForTest(0)
	domain.ChangeStampReadySecondsForTest(value.INTERVAL_SECONDS * 2)

	sh := sqlhelper.NewTestSqlHelper(testDB)
	svc := dkui.NewDkUiService(sh)
	ctx := context.Background()

	req := dkui.CreateWatchEventRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
		TrackID:   2,
		TalkID:    3,
		SlotID:    1000,
	}
	statReq := dkui.GetStatusRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
	}

	// first time
	err := svc.CreateWatchEvent(ctx, req)
	assert.Nil(t, err)

	resp, err := svc.GetStatus(ctx, statReq)
	assert.Nil(t, err)

	assert.Equal(t, int32(value.INTERVAL_SECONDS), resp.WatchedTalks.WatchingTime[1000])
	assert.Len(t, resp.StampChallenges, 0)

	// second time
	err = svc.CreateWatchEvent(ctx, req)
	assert.Nil(t, err)

	resp, err = svc.GetStatus(ctx, statReq)
	assert.Nil(t, err)

	assert.Equal(t, int32(value.INTERVAL_SECONDS*2), resp.WatchedTalks.WatchingTime[1000])
	assertStampStatus(t, resp, 1000, "ready")

	// stamp
	stampReq := dkui.StampOnlineRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
		SlotID:    1000,
	}
	err = svc.StampOnline(ctx, stampReq)
	assert.Nil(t, err)

	resp, err = svc.GetStatus(ctx, statReq)
	assert.Nil(t, err)

	assert.Equal(t, int32(value.INTERVAL_SECONDS*2), resp.WatchedTalks.WatchingTime[1000])
	assertStampStatus(t, resp, 1000, "stamped")
}

func TestDkUiServiceImpl_StampOnSite(t *testing.T) {
	sh := sqlhelper.NewTestSqlHelper(testDB)
	svc := dkui.NewDkUiService(sh)
	ctx := context.Background()

	req := dkui.StampOnSiteRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
		TrackID:   2,
		TalkID:    3,
		SlotID:    1001,
	}
	statReq := dkui.GetStatusRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
	}

	err := svc.StampOnSite(ctx, req)
	assert.Nil(t, err)

	resp, err := svc.GetStatus(ctx, statReq)
	assert.Nil(t, err)

	assert.Equal(t, int32(value.TALK_SECONDS), resp.WatchedTalks.WatchingTime[1001])
	assertStampStatus(t, resp, 1001, "stamped")
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}
