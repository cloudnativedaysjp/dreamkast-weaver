package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/repo"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/dbgen"
	"dreamkast-weaver/internal/stacktrace"
)

type DkUiRepoImpl struct {
	q *dbgen.Queries
}

func NewDkUiRepo(db dbgen.DBTX) repo.DkUiRepo {
	q := dbgen.New(db)
	return &DkUiRepoImpl{q}
}

var _ repo.DkUiRepo = (*DkUiRepoImpl)(nil)

func (r *DkUiRepoImpl) GetTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*dmodel.StampChallenges, error) {
	data, err := r.q.GetTrailmapStamps(ctx, dbgen.GetTrailmapStampsParams{
		ConferenceName: confName.String(),
		ProfileID:      profileID.Value(),
	})
	if err != nil {
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, stacktrace.With(fmt.Errorf("get stamp challenges: %w", err))
			}
		}
	}

	return stampChallengeConv.fromDB(data.Stamps)
}

func (r *DkUiRepoImpl) InsertViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID, ev *dmodel.ViewEvent) error {
	if err := r.q.InsertViewEvents(ctx, dbgen.InsertViewEventsParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
		TrackID:        ev.TrackID.Value(),
		TalkID:         ev.TalkID.Value(),
		SlotID:         ev.SlotID.Value(),
		ViewingSeconds: ev.ViewingSeconds.Value(),
	}); err != nil {
		return stacktrace.With(fmt.Errorf("insert view event: %w", err))
	}
	return nil
}

func (r *DkUiRepoImpl) ListViewEvents(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*dmodel.ViewEvents, error) {
	data, err := r.q.ListViewEvents(ctx, dbgen.ListViewEventsParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, stacktrace.With(fmt.Errorf("list view event: %w", err))
		}
	}

	return viewEventConv.fromDB(data)
}

func (r *DkUiRepoImpl) UpsertTrailMapStamps(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *dmodel.StampChallenges) error {
	buf, err := stampChallengeConv.toDB(scs)
	if err != nil {
		return stacktrace.With(err)
	}

	if err := r.q.UpsertTrailmapStamp(ctx, dbgen.UpsertTrailmapStampParams{
		ConferenceName: string(confName.Value()),
		ProfileID:      profileID.Value(),
		Stamps:         buf,
	}); err != nil {
		return stacktrace.With(fmt.Errorf("upsert stamp challenges: %w", err))
	}
	return nil
}

var stampChallengeConv _stampChallengeConv

type _stampChallengeConv struct{}

type _stampChallenge struct {
	SlotID    int32
	Condition string
	UpdatedAt time.Time
}

func (_stampChallengeConv) toDB(v *dmodel.StampChallenges) (json.RawMessage, error) {
	conv := func(dsc *dmodel.StampChallenge) *_stampChallenge {
		return &_stampChallenge{
			SlotID:    dsc.SlotID.Value(),
			Condition: string(dsc.Condition.Value()),
			UpdatedAt: dsc.UpdatedAt,
		}
	}

	var stamps []_stampChallenge
	for _, p := range v.Items {
		st := p
		stamps = append(stamps, *conv(&st))
	}

	buf, err := json.Marshal(stamps)
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("convert stamp challenges to DB: %w", err))
	}
	return json.RawMessage(buf), nil
}

func (_stampChallengeConv) fromDB(v json.RawMessage) (*dmodel.StampChallenges, error) {
	conv := func(sc *_stampChallenge) (*dmodel.StampChallenge, error) {
		slotID, err := value.NewSlotID(sc.SlotID)
		if err != nil {
			return nil, err
		}
		cond, err := value.NewStampCondition(value.StampConditionKind(sc.Condition))
		if err != nil {
			return nil, err
		}

		return &dmodel.StampChallenge{
			SlotID:    slotID,
			Condition: cond,
			UpdatedAt: sc.UpdatedAt,
		}, nil
	}

	if v == nil {
		return &dmodel.StampChallenges{}, nil
	}

	var stamps []_stampChallenge
	if err := json.Unmarshal(v, &stamps); err != nil {
		return nil, stacktrace.With(fmt.Errorf("convert stamp challenges from DB: %w", err))
	}

	var items []dmodel.StampChallenge
	for _, p := range stamps {
		st := p
		dst, err := conv(&st)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert stamp challenges from DB: %w", err))
		}
		items = append(items, *dst)
	}

	return &dmodel.StampChallenges{Items: items}, nil
}

var viewEventConv _viewEventConv

type _viewEventConv struct{}

func (_viewEventConv) fromDB(v []dbgen.ViewEvent) (*dmodel.ViewEvents, error) {
	conv := func(v *dbgen.ViewEvent) (*dmodel.ViewEvent, error) {
		trackID, err := value.NewTrackID(v.TrackID)
		if err != nil {
			return nil, err
		}
		talkID, err := value.NewTalkID(v.TalkID)
		if err != nil {
			return nil, err
		}
		slotID, err := value.NewSlotID(v.SlotID)
		if err != nil {
			return nil, err
		}
		viewingSeconds, err := value.NewViewingSeconds(v.ViewingSeconds)
		if err != nil {
			return nil, err
		}
		return &dmodel.ViewEvent{
			TrackID:        trackID,
			TalkID:         talkID,
			SlotID:         slotID,
			ViewingSeconds: viewingSeconds,
			CreatedAt:      v.CreatedAt,
		}, nil
	}

	var items []dmodel.ViewEvent

	for _, p := range v {
		ev := p
		dev, err := conv(&ev)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view event from DB: %w", err))
		}
		items = append(items, *dev)
	}

	return &dmodel.ViewEvents{Items: items}, nil
}

// func (_viewEventConv) toDB(confName value.ConfName, profileID value.ProfileID, v *dmodel.ViewEvents) []ViewEvent {
// 	conv := func(dev *dmodel.ViewEvent) *ViewEvent {
// 		return &ViewEvent{
// 			ConferenceName: string(confName.Value()),
// 			ProfileID:      profileID.Value(),
// 			TrackID:        dev.TrackID.Value(),
// 			TalkID:         dev.TalkID.Value(),
// 			SlotID:         dev.SlotID.Value(),
// 			ViewingSeconds: dev.ViewingSeconds.Value(),
// 			CreatedAt:      dev.CreatedAt,
// 		}
// 	}
//
// 	var events []ViewEvent
// 	for _, ev := range v.Items {
// 		events = append(events, *conv(&ev))
// 	}
//
// 	return events
// }

func (r *DkUiRepoImpl) InsertTrackViewer(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) error {
	if err := r.q.InsertTrackViewer(ctx, dbgen.InsertTrackViewerParams{
		ProfileID: profileID.Value(),
		TrackName: trackName.String(),
		TalkID:    talkID.Value(),
	}); err != nil {
		return stacktrace.With(fmt.Errorf("insert viewing track: %w", err))
	}
	return nil
}

func (r *DkUiRepoImpl) ListTrackViewer(ctx context.Context, from, to time.Time) (*dmodel.TrackViewers, error) {
	tvs, err := r.q.ListTrackViewer(ctx, dbgen.ListTrackViewerParams{
		FromCreatedAt: from,
		ToCreatedAt:   to,
	})
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("list viewing track: %w", err))
	}

	return trackViewerConv.fromDB(tvs)
}

var trackViewerConv _trackViewerConv

type _trackViewerConv struct{}

func (_trackViewerConv) fromDB(tvs []dbgen.TrackViewer) (*dmodel.TrackViewers, error) {
	conv := func(v *dbgen.TrackViewer) (*dmodel.TrackViewer, error) {
		tn, err := value.NewTrackName(v.TrackName)
		if err != nil {
			return nil, err
		}
		pID, err := value.NewProfileID(v.ProfileID)
		if err != nil {
			return nil, err
		}

		return &dmodel.TrackViewer{
			TrackName: tn,
			ProfileID: pID,
			CreatedAt: v.CreatedAt,
		}, nil
	}

	var items []dmodel.TrackViewer

	for _, tv := range tvs {
		v := tv
		dv, err := conv(&v)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert track viewer from DB: %w", err))
		}
		items = append(items, *dv)
	}

	return &dmodel.TrackViewers{Items: items}, nil
}

// func (_trackViewerConv) toDB(v *dmodel.TrackViewers) []TrackViewer {
// 	conv := func(dtv *dmodel.TrackViewer) *TrackViewer {
// 		return &TrackViewer{
// 			TrackName: dtv.TrackName.String(),
// 			ProfileID: dtv.ProfileID.Value(),
// 			CreatedAt: dtv.CreatedAt,
// 		}
// 	}
// 	var tvs []TrackViewer
// 	for _, vc := range v.Items {
// 		tvs = append(tvs, *conv(&vc))
// 	}
// 	return tvs
// }

type CfpRepoImpl struct {
	q *dbgen.Queries
}

func NewCfpRepo(db dbgen.DBTX) repo.CfpRepo {
	q := dbgen.New(db)
	return &CfpRepoImpl{q}
}

func (r *CfpRepoImpl) ListCfpVotes(ctx context.Context, confName value.ConfName, vt value.VotingTerm) (*dmodel.CfpVotes, error) {
	s, e := vt.Value()
	req := dbgen.ListCfpVotesParams{
		ConferenceName: string(confName.Value()),
		Start:          s,
		End:            e,
	}

	votes, err := r.q.ListCfpVotes(ctx, req)
	if err != nil {
		return nil, stacktrace.With(fmt.Errorf("list cfp vote: %w", err))
	}

	return cfpVoteConv.fromDB(votes)
}

func (r *CfpRepoImpl) InsertCfpVote(ctx context.Context, confName value.ConfName, talkID value.TalkID, clientIp net.IP) error {
	req := dbgen.InsertCfpVoteParams{
		ConferenceName: string(confName.Value()),
		TalkID:         talkID.Value(),
		ClientIp: sql.NullString{
			String: clientIp.String(),
			Valid:  true,
		},
	}

	if err := r.q.InsertCfpVote(ctx, req); err != nil {
		return stacktrace.With(fmt.Errorf("incert cfp vote: %w", err))
	}
	return nil
}

var cfpVoteConv _cfpVoteConv

type _cfpVoteConv struct{}

func (_cfpVoteConv) fromDB(v []dbgen.CfpVote) (*dmodel.CfpVotes, error) {
	conv := func(v *dbgen.CfpVote) (*dmodel.CfpVote, error) {
		talkID, err := value.NewTalkID(v.TalkID)
		if err != nil {
			return nil, err
		}
		ip := net.ParseIP(v.ClientIp.String)
		return &dmodel.CfpVote{
			TalkID:    talkID,
			ClientIp:  ip,
			CreatedAt: v.CreatedAt,
		}, nil
	}

	var items []dmodel.CfpVote
	for _, p := range v {
		cv := p
		dcv, err := conv(&cv)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view event from DB: %w", err))
		}
		items = append(items, *dcv)
	}

	return &dmodel.CfpVotes{Items: items}, nil
}
