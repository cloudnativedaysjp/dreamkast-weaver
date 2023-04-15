package dkui

import (
	"context"
	"database/sql"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/repo"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/graph/model"
	"dreamkast-weaver/internal/sqlhelper"
	"errors"
	"log"

	"github.com/ServiceWeaver/weaver"
)

type Service interface {
	CreateViewEvent(ctx context.Context, req model.CreateViewEventInput) error
	StampOnline(ctx context.Context, req model.StampOnlineInput) error
	StampOnSite(ctx context.Context, req model.StampOnSiteInput) error
	ViewingSlots(ctx context.Context, confName model.ConfName, profileID int) ([]*model.ViewingSlot, error)
	StampChallenges(ctx context.Context, confName model.ConfName, profileID int) ([]*model.StampChallenge, error)
}

type ServiceImpl struct {
	weaver.Implements[Service]
	weaver.WithConfig[config]

	sh     *sqlhelper.SqlHelper
	domain domain.DkUiDomain
}

var _ Service = (*ServiceImpl)(nil)

type config struct {
	DBUser     string `toml:"db_user"`
	DBPassword string `toml:"db_password"`
	DBEndpoint string `toml:"db_endpoint"`
	DBPort     string `toml:"db_port"`
	DBName     string `toml:"db_name"`
}

func (c *config) SqlOption() *sqlhelper.SqlOption {
	return &sqlhelper.SqlOption{
		User:     c.DBUser,
		Password: c.DBPassword,
		Endpoint: c.DBEndpoint,
		Port:     c.DBPort,
		DbName:   c.DBName,
	}
}

func NewDkUiService(sh *sqlhelper.SqlHelper) Service {
	return &ServiceImpl{sh: sh}
}

func (s *ServiceImpl) Init(ctx context.Context) error {
	cfg := s.Config()
	log.Printf("config: %#v\n", cfg)
	sh, err := sqlhelper.NewSqlHelper(cfg.SqlOption())
	if err != nil {
		return err
	}
	s.sh = sh
	return nil
}

func (v *ServiceImpl) CreateViewEvent(ctx context.Context, req model.CreateViewEventInput) error {
	r := repo.NewDkUiRepo(v.sh.DB())

	var mErr, err error
	confName, err := value.NewConfName(value.ConferenceKind(req.ConfName))
	errors.Join(mErr, err)
	profileID, err := value.NewProfileID(int32(req.ProfileID))
	errors.Join(mErr, err)
	trackID, err := value.NewTrackID(int32(req.TrackID))
	errors.Join(mErr, err)
	talkID, err := value.NewTalkID(int32(req.TalkID))
	errors.Join(mErr, err)
	slotID, err := value.NewSlotID(int32(req.SlotID))
	errors.Join(mErr, err)
	if err != nil {
		return err
	}

	devents, err := r.ListViewEvents(ctx, confName, profileID)
	if err != nil {
		return err
	}
	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return err
	}

	ev, err := v.domain.CreateOnlineViewEvent(trackID, talkID, slotID, dstamps, devents)
	if err != nil {
		return err
	}

	if err := v.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		r := repo.NewDkUiRepo(tx)
		if err := r.InsertViewEvents(ctx, confName, profileID, ev); err != nil {
			return err
		}
		if err := r.UpsertTrailMapStamps(ctx, confName, profileID, dstamps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (v *ServiceImpl) ViewingSlots(ctx context.Context, _confName model.ConfName, _profileID int) ([]*model.ViewingSlot, error) {
	r := repo.NewDkUiRepo(v.sh.DB())

	var mErr, err error
	confName, err := value.NewConfName(value.ConferenceKind(_confName.String()))
	errors.Join(mErr, err)
	profileID, err := value.NewProfileID(int32(_profileID))
	errors.Join(mErr, err)

	devents, err := r.ListViewEvents(ctx, confName, profileID)
	if err != nil {
		return nil, err
	}

	var viewingSlots []*model.ViewingSlot
	for k, v := range devents.ViewingSeconds() {
		viewingSlots = append(viewingSlots, &model.ViewingSlot{
			SlotID:      int(k.Value()),
			ViewingTime: int(v),
		})
	}

	return viewingSlots, nil
}

func (v *ServiceImpl) StampChallenges(ctx context.Context, _confName model.ConfName, _profileID int) ([]*model.StampChallenge, error) {
	r := repo.NewDkUiRepo(v.sh.DB())

	var mErr, err error
	confName, err := value.NewConfName(value.ConferenceKind(_confName.String()))
	errors.Join(mErr, err)
	profileID, err := value.NewProfileID(int32(_profileID))
	errors.Join(mErr, err)

	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return nil, err
	}

	var stamps []*model.StampChallenge
	for _, dst := range dstamps.Items {
		stamps = append(stamps, &model.StampChallenge{
			SlotID:    int(dst.SlotID.Value()),
			Condition: model.ChallengeCondition(dst.Condition.Value()),
			UpdatedAt: int(dst.UpdatedAt.Unix()),
		})
	}

	return stamps, nil
}

func (v *ServiceImpl) StampOnline(ctx context.Context, req model.StampOnlineInput) error {
	r := repo.NewDkUiRepo(v.sh.DB())

	var mErr, err error
	confName, err := value.NewConfName(value.ConferenceKind(req.ConfName))
	errors.Join(mErr, err)
	profileID, err := value.NewProfileID(int32(req.ProfileID))
	errors.Join(mErr, err)
	slotID, err := value.NewSlotID(int32(req.SlotID))
	errors.Join(mErr, err)
	if err != nil {
		return err
	}

	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return err
	}

	if err := v.domain.StampOnline(slotID, dstamps); err != nil {
		return err
	}

	if err := r.UpsertTrailMapStamps(ctx, confName, profileID, dstamps); err != nil {
		return err
	}

	return nil
}

func (v *ServiceImpl) StampOnSite(ctx context.Context, req model.StampOnSiteInput) error {
	r := repo.NewDkUiRepo(v.sh.DB())

	var mErr, err error
	confName, err := value.NewConfName(value.ConferenceKind(req.ConfName))
	errors.Join(mErr, err)
	profileID, err := value.NewProfileID(int32(req.ProfileID))
	errors.Join(mErr, err)
	trackID, err := value.NewTrackID(int32(req.TrackID))
	errors.Join(mErr, err)
	talkID, err := value.NewTalkID(int32(req.TalkID))
	errors.Join(mErr, err)
	slotID, err := value.NewSlotID(int32(req.SlotID))
	errors.Join(mErr, err)
	if err != nil {
		return err
	}

	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return err
	}

	ev, err := v.domain.StampOnSite(trackID, talkID, slotID, dstamps)
	if err != nil {
		return err
	}

	if err := v.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		r := repo.NewDkUiRepo(tx)
		if err := r.InsertViewEvents(ctx, confName, profileID, ev); err != nil {
			return err
		}
		if err := r.UpsertTrailMapStamps(ctx, confName, profileID, dstamps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
