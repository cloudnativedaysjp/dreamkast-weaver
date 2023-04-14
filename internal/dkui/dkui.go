package dkui

import (
	"context"
	"database/sql"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/repo"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/sqlhelper"
	"errors"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type DkUiService interface {
	CreateWatchEvent(ctx context.Context, req CreateWatchEventRequest) error
	StampOnline(ctx context.Context, req StampOnlineRequest) error
	StampOnSite(ctx context.Context, req StampOnSiteRequest) error
	GetStatus(ctx context.Context, req GetStatusRequest) (*StatusResponse, error)
}

type CreateWatchEventRequest struct {
	weaver.AutoMarshal
	ConfName  string
	ProfileID int
	TrackID   int
	TalkID    int
	SlotID    int
}

type StampOnlineRequest struct {
	ConfName  string
	ProfileID int
	SlotID    int
}

type StampOnSiteRequest struct {
	ConfName  string
	ProfileID int
	TrackID   int
	TalkID    int
	SlotID    int
}

type GetStatusRequest struct {
	ConfName  string
	ProfileID int
}

type StatusResponse struct {
	WatchedTalks
	StampChallenges []StampChallenge
}

type WatchedTalks struct {
	WatchingTime  map[int32]int32
	PrevTimestamp int
}

type StampChallenge struct {
	SlotID    int
	Condition string
	UpdatedAt time.Time
}

type DkUiServiceImpl struct {
	weaver.Implements[DkUiService]
	weaver.WithConfig[config]

	sh     *sqlhelper.SqlHelper
	domain domain.DkUiDomain
}

var _ DkUiService = (*DkUiServiceImpl)(nil)

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

func NewDkUiService(sh *sqlhelper.SqlHelper) DkUiService {
	return &DkUiServiceImpl{sh: sh}
}

func (v *DkUiServiceImpl) CreateWatchEvent(ctx context.Context, req CreateWatchEventRequest) error {
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

	devents, err := r.ListWatchEvents(ctx, confName, profileID)
	if err != nil {
		return err
	}
	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return err
	}

	ev, err := v.domain.CreateOnlineWatchEvent(trackID, talkID, slotID, dstamps, devents)
	if err != nil {
		return err
	}

	if err := v.sh.RunTX(ctx, func(ctx context.Context, tx *sql.Tx) error {
		r := repo.NewDkUiRepo(tx)
		if err := r.InsertWatchEvents(ctx, confName, profileID, ev); err != nil {
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

func (v *DkUiServiceImpl) GetStatus(ctx context.Context, req GetStatusRequest) (*StatusResponse, error) {
	r := repo.NewDkUiRepo(v.sh.DB())

	var mErr, err error
	confName, err := value.NewConfName(value.ConferenceKind(req.ConfName))
	errors.Join(mErr, err)
	profileID, err := value.NewProfileID(int32(req.ProfileID))
	errors.Join(mErr, err)

	devents, err := r.ListWatchEvents(ctx, confName, profileID)
	if err != nil {
		return nil, err
	}
	dstamps, err := r.GetTrailMapStamps(ctx, confName, profileID)
	if err != nil {
		return nil, err
	}

	viewingSeconds := map[int32]int32{}
	for k, v := range devents.ViewingSeconds() {
		viewingSeconds[k.Value()] = v
	}
	var stamps []StampChallenge
	for _, dst := range dstamps.Items {
		stamps = append(stamps, StampChallenge{
			SlotID:    int(dst.SlotID.Value()),
			Condition: string(dst.Condition.Value()),
			UpdatedAt: dst.UpdatedAt,
		})
	}

	return &StatusResponse{
		WatchedTalks: WatchedTalks{
			WatchingTime:  viewingSeconds,
			PrevTimestamp: int(devents.LastCreated().Unix()),
		},
		StampChallenges: stamps,
	}, nil
}

func (v *DkUiServiceImpl) StampOnline(ctx context.Context, req StampOnlineRequest) error {
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

func (v *DkUiServiceImpl) StampOnSite(ctx context.Context, req StampOnSiteRequest) error {
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
		if err := r.InsertWatchEvents(ctx, confName, profileID, ev); err != nil {
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
