package application

import (
	"context"
	"database/sql"

	derrors "dreamkast-weaver/internal/domain/errors"
	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/repo"
	"dreamkast-weaver/internal/pkg/logger"
	"dreamkast-weaver/internal/pkg/sqlhelper"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type StampRallyApp interface {
	CreateViewEvent(ctx context.Context, profile Profile, req CreateViewEventRequest) error
	StampOnline(ctx context.Context, profile Profile, slotID value.SlotID) error
	StampOnSite(ctx context.Context, profile Profile, req StampRequest) error
	ViewingEvents(ctx context.Context, profile Profile) (*dmodel.ViewEvents, error)
	StampChallenges(ctx context.Context, profile Profile) (*dmodel.StampChallenges, error)
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

type StampRallyAppImpl struct {
	sh         *sqlhelper.SqlHelper
	dkUiDomain dmodel.DkUiDomain
}

var _ StampRallyApp = (*StampRallyAppImpl)(nil)

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

func NewStampRallyApp(sh *sqlhelper.SqlHelper) StampRallyApp {
	return &StampRallyAppImpl{
		sh: sh,
	}
}

func (s *StampRallyAppImpl) handleError(ctx context.Context, msg string, err error) {
	logger := logger.FromCtx(ctx)
	if err != nil {
		if derrors.IsUserError(err) {
			logger.With("errorType", "user-side").Info(msg, err)
		} else {
			logger.With("stacktrace", stacktrace.Get(err)).Error(msg, err)
		}
	}
}

func (s *StampRallyAppImpl) CreateViewEvent(ctx context.Context, profile Profile, req CreateViewEventRequest) (err error) {
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

func (v *StampRallyAppImpl) ViewingEvents(ctx context.Context, profile Profile) (resp *dmodel.ViewEvents, err error) {
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

func (v *StampRallyAppImpl) StampChallenges(ctx context.Context, profile Profile) (resp *dmodel.StampChallenges, err error) {
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

func (v *StampRallyAppImpl) StampOnline(ctx context.Context, profile Profile, slotID value.SlotID) (err error) {
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

func (v *StampRallyAppImpl) StampOnSite(ctx context.Context, profile Profile, req StampRequest) (err error) {
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
