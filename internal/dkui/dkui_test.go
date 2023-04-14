package dkui_test

import (
	"context"
	"dreamkast-weaver/internal/dkui"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/repo"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/sqlhelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDkUiServiceImpl_CreateWatchEvent(t *testing.T) {
	domain.ChangeGuardSecondsForTest(0)

	sh := sqlhelper.NewTestSqlHelper("dkui")
	svc := dkui.NewDkUiService(sh)
	ctx := context.Background()

	req := dkui.CreateWatchEventRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
		TrackID:   2,
		TalkID:    3,
		SlotID:    1000,
	}

	err := svc.CreateWatchEvent(ctx, req)
	assert.Nil(t, err)

	// TODO check record
}

func TestDkUiServiceImpl_GetStatus(t *testing.T) {
	sh := sqlhelper.NewTestSqlHelper("dkui")
	svc := dkui.NewDkUiService(sh)
	ctx := context.Background()

	req := dkui.GetStatusRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
	}

	resp, err := svc.GetStatus(ctx, req)
	assert.Nil(t, err)

	assert.Greater(t, len(resp.WatchedTalks.WatchingTime), 0)
	assert.Greater(t, len(resp.StampChallenges), 0)
}

func TestDkUiServiceImpl_StampOnline(t *testing.T) {
	sh := sqlhelper.NewTestSqlHelper("dkui")
	svc := dkui.NewDkUiService(sh)
	ctx := context.Background()

	req := dkui.StampOnlineRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
		SlotID:    1000,
	}

	err := svc.StampOnline(ctx, req)
	assert.Nil(t, err)

	r := repo.NewDkUiRepo(sh.DB())

	resp, err := r.GetTrailMapStamps(ctx, value.CNDF2023, value.ProfileID(newProfileID(1)))
	assert.Nil(t, err)

	got := resp.Get(newSlotID(1000))
	assert.NotNil(t, got)
	assert.Equal(t, value.StampStamped, got.Condition)
}

func TestDkUiServiceImpl_StampOnSite(t *testing.T) {
	sh := sqlhelper.NewTestSqlHelper("dkui")
	svc := dkui.NewDkUiService(sh)
	ctx := context.Background()

	req := dkui.StampOnSiteRequest{
		ConfName:  "cndf2023",
		ProfileID: 1,
		TrackID:   2,
		TalkID:    3,
		SlotID:    1001,
	}

	err := svc.StampOnSite(ctx, req)
	assert.Nil(t, err)

	r := repo.NewDkUiRepo(sh.DB())

	resp, err := r.GetTrailMapStamps(ctx, value.CNDF2023, value.ProfileID(newProfileID(1)))
	assert.Nil(t, err)

	got := resp.Get(newSlotID(1001))
	assert.NotNil(t, got)
	assert.Equal(t, value.StampStamped, got.Condition)
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func newProfileID(v int32) value.ProfileID {
	id, err := value.NewProfileID(v)
	mustNil(err)
	return id
}

func newSlotID(v int32) value.SlotID {
	id, err := value.NewSlotID(v)
	mustNil(err)
	return id
}
