package dkui_test

import (
	"context"
	"dreamkast-weaver/internal/dkui"
	"dreamkast-weaver/internal/dkui/domain"
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
}
