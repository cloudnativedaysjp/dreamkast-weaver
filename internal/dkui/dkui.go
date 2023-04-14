package dkui

import (
	"context"
	"dreamkast-weaver/internal/dkui/repo"
	"dreamkast-weaver/internal/sqlhelper"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type DkUiService interface {
	CreateWatchEvent(ctx context.Context, req CreateWatchEventRequest) error
	StampOnline(ctx context.Context, req StampOnlineRequest) error
	StampOnSite(ctx context.Context, req StampOnSiteRequest) error
	GetStatus(ctx context.Context, req GetStatusRequest)
}

type CreateWatchEventRequest struct {
	weaver.AutoMarshal
	ConfName  string
	ProfileID int32
	TrackID   int32
	TalkID    int32
	SlotID    int32
}

type StampOnlineRequest struct {
	SlotID int32
}

type StampOnSiteRequest struct {
	TrackID int32
	TalkID  int32
	SlotID  int32
}

type GetStatusRequest struct {
	ConfName  string
	ProfileID int32
}

type StatusResponse struct {
	WatchedTalks
	StampChallenges []StampChallenge
}

type WatchedTalks struct {
	WatchingTime  map[int32]int
	PrevTimestamp int
}

type StampChallenge struct {
	SlotID    int32
	Condition string
	UpdatedAt time.Time
}

type DkUiServiceImpl struct {
	weaver.Implements[DkUiService]
	weaver.WithConfig[config]

	sh *sqlhelper.SqlHelper
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
	r := repo.New(v.sh.DB())

	if err := r.InsertWatchEvents(ctx, repo.InsertWatchEventsParams{
		ProfileID:      req.ProfileID,
		ConferenceName: req.ConfName,
		TrackID:        req.TrackID,
		TalkID:         req.TalkID,
		ViewingSeconds: 120,
	}); err != nil {
		return fmt.Errorf("incert viewing event: %w", err)
	}

	return nil
}

func (v *DkUiServiceImpl) GetStatus(ctx context.Context, req GetStatusRequest) {
	// r := repo.new(v.sh.db())

	// events, err := r.ListWatchEvents(ctx, repo.ListWatchEventsParams{
	// 	ConferenceName: req.ConfName,
	// 	ProfileID:      req.ProfileID,
	// })

	// stamps, err := r.GetTrailmapStamps(ctx, repo.GetTrailmapStampsParams{
	// 	ConferenceName: "",
	// 	ProfileID:      0,
	// })

	panic("unimplemented")
}

func (v *DkUiServiceImpl) StampOnSite(ctx context.Context, req StampOnSiteRequest) error {
	panic("unimplemented")
}

func (v *DkUiServiceImpl) StampOnline(ctx context.Context, req StampOnlineRequest) error {
	panic("unimplemented")
}
