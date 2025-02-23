package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"golang.org/x/exp/slog"

	derrors "dreamkast-weaver/internal/domain/errors"
	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/repo"
	"dreamkast-weaver/internal/pkg/logger"
	"dreamkast-weaver/internal/pkg/sqlhelper"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

const (
	envPushGatewayEndpoint = "PROM_PUSHGATEWAY_ENDPOINT"
)

var (
	viewerCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "dkw",
		Subsystem: "dkui",
		Name:      "viewer_count",
		Help:      "Number of viewer_count",
	}, []string{"trackName"})
)

type DkUiService interface {
	CreateViewEvent(ctx context.Context, profile Profile, req CreateViewEventRequest) error
	StampOnline(ctx context.Context, profile Profile, slotID value.SlotID) error
	StampOnSite(ctx context.Context, profile Profile, req StampRequest) error
	ViewingEvents(ctx context.Context, profile Profile) (*dmodel.ViewEvents, error)
	StampChallenges(ctx context.Context, profile Profile) (*dmodel.StampChallenges, error)
	ListViewerCounts(ctx context.Context, useCache bool) (*dmodel.ViewerCounts, error)
	ViewTrack(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) error
}

type Profile struct {
	ID       value.ProfileID
	ConfName value.ConfName
}

type CreateViewEventRequest struct {
	TrackID value.TrackID
	TalkID  value.TalkID
	SlotID  value.SlotID
}

type StampRequest struct {
	TrackID value.TrackID
	TalkID  value.TalkID
	SlotID  value.SlotID
}

type DkUiServiceImpl struct {
	sh         *sqlhelper.SqlHelper
	pusher     *push.Pusher
	dkUiDomain dmodel.DkUiDomain
	cache      dmodel.ViewerCounts
}

var _ DkUiService = (*DkUiServiceImpl)(nil)

type Config struct {
	DBUser              string `toml:"db_user"`
	DBPassword          string `toml:"db_password"`
	DBEndpoint          string `toml:"db_endpoint"`
	DBPort              string `toml:"db_port"`
	DBName              string `toml:"db_name"`
	PushGatewayEndpoint string `toml:"push_gateway_endpoint"`
}

func (c *Config) SqlOption() *sqlhelper.SqlOption {
	return &sqlhelper.SqlOption{
		User:     c.DBUser,
		Password: c.DBPassword,
		Endpoint: c.DBEndpoint,
		Port:     c.DBPort,
		DbName:   c.DBName,
	}
}

func NewDkUiService(sh *sqlhelper.SqlHelper) DkUiService {
	return &DkUiServiceImpl{
		sh: sh,
	}
}

func (s *DkUiServiceImpl) handleError(ctx context.Context, msg string, err error) {
	logger := logger.FromCtx(ctx)
	if err != nil {
		if derrors.IsUserError(err) {
			logger.With("errorType", "user-side").Info(msg, err)
		} else {
			logger.With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *DkUiServiceImpl) CreateViewEvent(ctx context.Context, profile Profile, req CreateViewEventRequest) (err error) {
	defer func() {
		s.handleError(ctx, "create viewEvent", err)
	}()

	vr := repo.NewViewEventRepo(s.sh.DB())
	tr := repo.NewTrailMapStampsRepo(s.sh.DB())

	devents, err := vr.List(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}
	dstamps, err := tr.Get(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}

	ev, err := s.dkUiDomain.CreateOnlineViewEvent(req.TrackID, req.TalkID, req.SlotID, dstamps, devents)
	if err != nil {
		return err
	}

	if err := s.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		ver := repo.NewViewEventRepo(tx)
		tsr := repo.NewTrailMapStampsRepo(s.sh.DB())

		if err := ver.Insert(ctx, profile.ConfName, profile.ID, ev); err != nil {
			return err
		}
		if err := tsr.Upsert(ctx, profile.ConfName, profile.ID, dstamps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (v *DkUiServiceImpl) ViewingEvents(ctx context.Context, profile Profile) (resp *dmodel.ViewEvents, err error) {
	defer func() {
		v.handleError(ctx, "get viewingEvents", err)
	}()

	r := repo.NewViewEventRepo(v.sh.DB())

	resp, err = r.List(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *DkUiServiceImpl) StampChallenges(ctx context.Context, profile Profile) (resp *dmodel.StampChallenges, err error) {
	defer func() {
		v.handleError(ctx, "get stampChallenges", err)
	}()

	r := repo.NewTrailMapStampsRepo(v.sh.DB())

	resp, err = r.Get(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *DkUiServiceImpl) StampOnline(ctx context.Context, profile Profile, slotID value.SlotID) (err error) {
	defer func() {
		v.handleError(ctx, "stamp from online", err)
	}()

	r := repo.NewTrailMapStampsRepo(v.sh.DB())

	dstamps, err := r.Get(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}

	if err := v.dkUiDomain.StampOnline(slotID, dstamps); err != nil {
		return err
	}

	if err := r.Upsert(ctx, profile.ConfName, profile.ID, dstamps); err != nil {
		return err
	}

	return nil
}

func (v *DkUiServiceImpl) StampOnSite(ctx context.Context, profile Profile, req StampRequest) (err error) {
	defer func() {
		v.handleError(ctx, "stamp from onsite", err)
	}()

	r := repo.NewTrailMapStampsRepo(v.sh.DB())

	dstamps, err := r.Get(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}

	ev, err := v.dkUiDomain.StampOnSite(req.TrackID, req.TalkID, req.SlotID, dstamps)
	if err != nil {
		return err
	}

	if err := v.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err := repo.NewViewEventRepo(tx).Insert(ctx, profile.ConfName, profile.ID, ev); err != nil {
			return err
		}
		if err := repo.NewTrailMapStampsRepo(tx).Upsert(ctx, profile.ConfName, profile.ID, dstamps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *DkUiServiceImpl) ListViewerCounts(ctx context.Context, useCache bool) (dvc *dmodel.ViewerCounts, err error) {
	defer func() {
		s.handleError(ctx, "list viewer count", err)
	}()
	if !useCache {
		if _, err := s.getViewerCount(ctx); err != nil {
			return nil, err
		}
	}
	return &s.cache, nil
}

func (s *DkUiServiceImpl) ViewTrack(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) (err error) {
	defer func() {
		s.handleError(ctx, "viewing track", err)
	}()

	r := repo.NewTrackViewerRepo(s.sh.DB())
	if err := r.Insert(ctx, profileID, trackName, talkID); err != nil {
		return err
	}
	return nil
}

func (s *DkUiServiceImpl) measureViewerCount(ctx context.Context) {
	logger := logger.FromCtx(ctx)
	go func() {
		for {
			vc, err := s.getViewerCount(ctx)
			if err != nil {
				logger.Warn("failed push metrics", slog.String("err", err.Error()))
			}

			for _, v := range vc.Items {
				viewerCount.WithLabelValues(v.TrackName.String()).Set(float64(v.Count))
			}

			if err := s.pusher.Push(); err != nil {
				logger.Warn("failed push metrics", slog.String("err", err.Error()))
			}

			time.Sleep(value.METRICS_UPDATE_INTERVAL * time.Second)
		}
	}()
}

func (s *DkUiServiceImpl) getViewerCount(ctx context.Context) (*dmodel.ViewerCounts, error) {
	to := time.Now().UTC()
	from := to.Add(-1 * value.TIMEWINDOW_VIEWER_COUNT * time.Second)

	r := repo.NewTrackViewerRepo(s.sh.DB())
	dtv, err := r.List(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dvc := dtv.GetViewerCounts()
	s.cache = dvc

	return &dvc, nil
}
