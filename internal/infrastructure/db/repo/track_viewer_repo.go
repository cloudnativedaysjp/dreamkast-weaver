package repo

import (
	"context"
	"fmt"
	"time"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/repo"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/dbgen"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type TrackViewerRepoImpl struct {
	q *dbgen.Queries
}

func NewTrackViewerRepo(db dbgen.DBTX) repo.TrackViewerRepo {
	q := dbgen.New(db)
	return &TrackViewerRepoImpl{q}
}

var _ repo.TrackViewerRepo = (*TrackViewerRepoImpl)(nil)

func (r *TrackViewerRepoImpl) Insert(ctx context.Context, profileID value.ProfileID, trackName value.TrackName, talkID value.TalkID) error {
	if err := r.q.InsertTrackViewer(ctx, dbgen.InsertTrackViewerParams{
		ProfileID: profileID.Value(),
		TrackName: trackName.String(),
		TalkID:    talkID.Value(),
	}); err != nil {
		return stacktrace.With(fmt.Errorf("insert viewing track: %w", err))
	}
	return nil
}

func (r *TrackViewerRepoImpl) List(ctx context.Context, from, to time.Time) (*dmodel.TrackViewers, error) {
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
