package dkui

import (
	"context"
	"database/sql"
	"os"
	"time"

	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/repo"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/stacktrace"

	"github.com/ServiceWeaver/weaver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"golang.org/x/exp/slog"
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

type Service interface {
	CreateViewEvent(ctx context.Context, profile Profile, req CreateViewEventRequest) error
	StampOnline(ctx context.Context, profile Profile, slotID value.SlotID) error
	StampOnSite(ctx context.Context, profile Profile, req StampRequest) error
	ViewingEvents(ctx context.Context, profile Profile) (*domain.ViewEvents, error)
	StampChallenges(ctx context.Context, profile Profile) (*domain.StampChallenges, error)
	ListViewerCounts(ctx context.Context, useCache bool) (*domain.ViewerCounts, error)
	ViewTrack(ctx context.Context, profileID value.ProfileID, trackName value.TrackName) error
}

type Profile struct {
	weaver.AutoMarshal
	ID       value.ProfileID
	ConfName value.ConfName
}

type CreateViewEventRequest struct {
	weaver.AutoMarshal
	TrackID value.TrackID
	TalkID  value.TalkID
	SlotID  value.SlotID
}

type StampRequest struct {
	weaver.AutoMarshal
	TrackID value.TrackID
	TalkID  value.TalkID
	SlotID  value.SlotID
}

type ServiceImpl struct {
	weaver.Implements[Service]
	weaver.WithConfig[Config]

	sh     *sqlhelper.SqlHelper
	pusher *push.Pusher
	domain domain.DkUiDomain
	cache  domain.ViewerCounts
}

var _ Service = (*ServiceImpl)(nil)

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

func NewService(sh *sqlhelper.SqlHelper) Service {
	return &ServiceImpl{
		sh: sh,
	}
}

func (s *ServiceImpl) Init(ctx context.Context) error {
	opt := s.Config().SqlOption()
	if err := opt.Validate(); err != nil {
		opt = sqlhelper.NewOptionFromEnv("dkui")
	}
	sh, err := sqlhelper.NewSqlHelper(opt)
	if err != nil {
		return err
	}
	s.sh = sh

	endpoint := s.Config().PushGatewayEndpoint
	if endpoint == "" {
		endpoint = os.Getenv(envPushGatewayEndpoint)
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(viewerCount)
	s.pusher = push.New(endpoint, "dkw_dkui").Gatherer(registry)

	s.measureViewerCount(ctx)
	return nil
}

func (s *ServiceImpl) HandleError(msg string, err error) {
	if err != nil {
		if derrors.IsUserError(err) {
			s.Logger().With("errorType", "user-side").Info(msg, err)
		} else {
			s.Logger().With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *ServiceImpl) CreateViewEvent(ctx context.Context, profile Profile, req CreateViewEventRequest) (err error) {
	defer func() {
		s.HandleError("create viewEvent", err)
	}()

	r := repo.NewDkUiRepo(s.sh.DB())

	devents, err := r.ListViewEvents(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}
	dstamps, err := r.GetTrailMapStamps(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}

	ev, err := s.domain.CreateOnlineViewEvent(req.TrackID, req.TalkID, req.SlotID, dstamps, devents)
	if err != nil {
		return err
	}

	if err := s.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		r := repo.NewDkUiRepo(tx)
		if err := r.InsertViewEvents(ctx, profile.ConfName, profile.ID, ev); err != nil {
			return err
		}
		if err := r.UpsertTrailMapStamps(ctx, profile.ConfName, profile.ID, dstamps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (v *ServiceImpl) ViewingEvents(ctx context.Context, profile Profile) (resp *domain.ViewEvents, err error) {
	defer func() {
		v.HandleError("get viewingEvents", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	resp, err = r.ListViewEvents(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *ServiceImpl) StampChallenges(ctx context.Context, profile Profile) (resp *domain.StampChallenges, err error) {
	defer func() {
		v.HandleError("get stampChallenges", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	resp, err = r.GetTrailMapStamps(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *ServiceImpl) StampOnline(ctx context.Context, profile Profile, slotID value.SlotID) (err error) {
	defer func() {
		v.HandleError("stamp from online", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	dstamps, err := r.GetTrailMapStamps(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}

	if err := v.domain.StampOnline(slotID, dstamps); err != nil {
		return err
	}

	if err := r.UpsertTrailMapStamps(ctx, profile.ConfName, profile.ID, dstamps); err != nil {
		return err
	}

	return nil
}

func (v *ServiceImpl) StampOnSite(ctx context.Context, profile Profile, req StampRequest) (err error) {
	defer func() {
		v.HandleError("stamp from onsite", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	dstamps, err := r.GetTrailMapStamps(ctx, profile.ConfName, profile.ID)
	if err != nil {
		return err
	}

	ev, err := v.domain.StampOnSite(req.TrackID, req.TalkID, req.SlotID, dstamps)
	if err != nil {
		return err
	}

	if err := v.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		r := repo.NewDkUiRepo(tx)
		if err := r.InsertViewEvents(ctx, profile.ConfName, profile.ID, ev); err != nil {
			return err
		}
		if err := r.UpsertTrailMapStamps(ctx, profile.ConfName, profile.ID, dstamps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) ListViewerCounts(ctx context.Context, useCache bool) (dvc *domain.ViewerCounts, err error) {
	defer func() {
		s.HandleError("list viewer count", err)
	}()
	if !useCache {
		if _, err := s.getViewerCount(ctx); err != nil {
			return nil, err
		}
	}
	return &s.cache, nil
}

func (s *ServiceImpl) ViewTrack(ctx context.Context, profileID value.ProfileID, trackName value.TrackName) (err error) {
	defer func() {
		s.HandleError("viewing track", err)
	}()

	r := repo.NewDkUiRepo(s.sh.DB())
	if err := r.InsertTrackViewer(ctx, profileID, trackName); err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) measureViewerCount(ctx context.Context) {
	logger := s.Logger()
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

func (s *ServiceImpl) getViewerCount(ctx context.Context) (*domain.ViewerCounts, error) {
	to := time.Now().UTC()
	from := to.Add(-1 * value.TIMEWINDOW_VIEWER_COUNT * time.Second)

	r := repo.NewDkUiRepo(s.sh.DB())
	dtv, err := r.ListTrackViewer(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dvc := dtv.GetViewerCounts()
	s.cache = dvc

	return &dvc, nil
}
