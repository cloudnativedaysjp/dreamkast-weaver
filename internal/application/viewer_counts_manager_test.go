package application

import (
	"context"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/pkg/sqlhelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewerCountsManagerImpl_ListTrackViewer(t *testing.T) {
	ctx := context.Background()
	sq, _ := sqlhelper.NewSqlHelper(&sqlhelper.SqlOption{
		User:     "user",
		Password: "password",
		Endpoint: "127.0.0.1",
		Port:     "13306",
		DbName:   "test_dkui",
	})

	svc := NewViewerCountManager(sq)

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

	dvc, err := svc.ListViewerCounts(ctx)
	assert.NoError(t, err)

	for _, v := range dvc.Items {
		assert.Equal(t, ans[v.TrackName], v.Count)
	}
}
