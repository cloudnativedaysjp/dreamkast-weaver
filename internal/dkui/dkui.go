package dkui

import (
	"context"
	"database/sql"
	"errors"

	"dreamkast-weaver/internal/derrors"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/repo"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/graph/model"
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/stacktrace"

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

func NewService(sh *sqlhelper.SqlHelper) Service {
	return &ServiceImpl{sh: sh}
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
	return nil
}

func (s *ServiceImpl) HandleError(msg string, err error) {
	if err != nil && !derrors.IsUserError(err) {
		s.Logger().With("stacktrace", stacktrace.Get(err)).Error(msg, err)
	}
}

func (s *ServiceImpl) CreateViewEvent(ctx context.Context, req model.CreateViewEventInput) (err error) {
	defer func() {
		s.HandleError("create viewEvent", err)
	}()

	r := repo.NewDkUiRepo(s.sh.DB())

	var e error
	confName, e := value.NewConfName(value.ConferenceKind(req.ConfName))
	err = errors.Join(err, e)
	profileID, e := value.NewProfileID(int32(req.ProfileID))
	err = errors.Join(err, e)
	trackID, e := value.NewTrackID(int32(req.TrackID))
	err = errors.Join(err, e)
	talkID, e := value.NewTalkID(int32(req.TalkID))
	err = errors.Join(err, e)
	slotID, e := value.NewSlotID(int32(req.SlotID))
	err = errors.Join(err, e)
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

	ev, err := s.domain.CreateOnlineViewEvent(trackID, talkID, slotID, dstamps, devents)
	if err != nil {
		return err
	}

	if err := s.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
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

func (v *ServiceImpl) ViewingSlots(ctx context.Context, _confName model.ConfName, _profileID int) (viewingSlots []*model.ViewingSlot, err error) {
	defer func() {
		v.HandleError("get viewingSlots", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	var e error
	confName, e := value.NewConfName(value.ConferenceKind(_confName.String()))
	err = errors.Join(err, e)
	profileID, e := value.NewProfileID(int32(_profileID))
	err = errors.Join(err, e)
	if err != nil {
		return nil, stacktrace.With(err)
	}

	devents, err := r.ListViewEvents(ctx, confName, profileID)
	if err != nil {
		return nil, err
	}

	for k, v := range devents.ViewingSeconds() {
		viewingSlots = append(viewingSlots, &model.ViewingSlot{
			SlotID:      int(k.Value()),
			ViewingTime: int(v),
		})
	}

	return viewingSlots, nil
}

func (v *ServiceImpl) StampChallenges(ctx context.Context, _confName model.ConfName, _profileID int) (stamps []*model.StampChallenge, err error) {
	defer func() {
		v.HandleError("get stampChallenges", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	var e error
	confName, e := value.NewConfName(value.ConferenceKind(_confName.String()))
	err = errors.Join(err, e)
	profileID, e := value.NewProfileID(int32(_profileID))
	err = errors.Join(err, e)
	if err != nil {
		return nil, err
	}

	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return nil, err
	}

	for _, dst := range dstamps.Items {
		stamps = append(stamps, &model.StampChallenge{
			SlotID:    int(dst.SlotID.Value()),
			Condition: model.ChallengeCondition(dst.Condition.Value()),
			UpdatedAt: int(dst.UpdatedAt.Unix()),
		})
	}

	return stamps, nil
}

func (v *ServiceImpl) StampOnline(ctx context.Context, req model.StampOnlineInput) (err error) {
	defer func() {
		v.HandleError("stamp from online", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	confName, e := value.NewConfName(value.ConferenceKind(req.ConfName))
	err = errors.Join(err, e)
	profileID, e := value.NewProfileID(int32(req.ProfileID))
	err = errors.Join(err, e)
	slotID, e := value.NewSlotID(int32(req.SlotID))
	err = errors.Join(err, e)
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

func (v *ServiceImpl) StampOnSite(ctx context.Context, req model.StampOnSiteInput) (err error) {
	defer func() {
		v.HandleError("stamp from onsite", err)
	}()

	r := repo.NewDkUiRepo(v.sh.DB())

	confName, e := value.NewConfName(value.ConferenceKind(req.ConfName))
	err = errors.Join(err, e)
	profileID, e := value.NewProfileID(int32(req.ProfileID))
	err = errors.Join(err, e)
	trackID, e := value.NewTrackID(int32(req.TrackID))
	err = errors.Join(err, e)
	talkID, e := value.NewTalkID(int32(req.TalkID))
	err = errors.Join(err, e)
	slotID, e := value.NewSlotID(int32(req.SlotID))
	err = errors.Join(err, e)
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
