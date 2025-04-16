package application

import (
	"context"
	"log/slog"
	"time"

	derrors "dreamkast-weaver/internal/domain/errors"
	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/repo"
	"dreamkast-weaver/internal/pkg/logger"
	"dreamkast-weaver/internal/pkg/metrics"
	"dreamkast-weaver/internal/pkg/sqlhelper"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type ViewerCountManager interface {
	Run(ctx context.Context)
	ListViewerCounts(ctx context.Context) (*dmodel.ViewerCounts, error)
	ViewTrack(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) error
}

type ViewerCountManagerImpl struct {
	sh      *sqlhelper.SqlHelper
	current dmodel.ViewerCounts
}

var _ ViewerCountManager = (*ViewerCountManagerImpl)(nil)

func NewViewerCountManager(sh *sqlhelper.SqlHelper) ViewerCountManager {
	return &ViewerCountManagerImpl{
		sh: sh,
	}
}

func (s *ViewerCountManagerImpl) Run(ctx context.Context) {
	ticker := time.NewTicker(value.METRICS_UPDATE_INTERVAL * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.measureViewerCount(ctx)
		}
	}
}

func (s *ViewerCountManagerImpl) ListViewerCounts(ctx context.Context) (dvc *dmodel.ViewerCounts, err error) {
	return &s.current, nil
}

func (s *ViewerCountManagerImpl) ViewTrack(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) (err error) {
	defer func() {
		s.handleError(ctx, "viewing track", err)
	}()

	r := repo.NewTrackViewerRepo(s.sh.DB())
	if err := r.Insert(ctx, profileID, trackName, talkID); err != nil {
		return err
	}
	return nil
}

func (s *ViewerCountManagerImpl) handleError(ctx context.Context, msg string, err error) {
	logger := logger.FromCtx(ctx)
	if err != nil {
		if derrors.IsUserError(err) {
			logger.With("errorType", "user-side").Info(msg, err)
		} else {
			logger.With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *ViewerCountManagerImpl) measureViewerCount(ctx context.Context) {
	logger := logger.FromCtx(ctx)

	vc, err := s.getViewerCount(ctx)
	if err != nil {
		logger.Warn("failed get metrics", slog.String("err", err.Error()))
	}

	for _, v := range vc.Items {
		metrics.ViewerCount.WithLabelValues(v.TrackName.String()).Set(float64(v.Count))
	}
}

func (s *ViewerCountManagerImpl) getViewerCount(ctx context.Context) (*dmodel.ViewerCounts, error) {
	to := time.Now().UTC()
	from := to.Add(-1 * value.TIMEWINDOW_VIEWER_COUNT * time.Second)

	r := repo.NewTrackViewerRepo(s.sh.DB())
	dtv, err := r.List(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dvc := dtv.GetViewerCounts()
	s.current = dvc

	return &dvc, nil
}
