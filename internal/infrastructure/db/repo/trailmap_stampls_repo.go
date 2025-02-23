package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dmodel "dreamkast-weaver/internal/domain/model"
	"dreamkast-weaver/internal/domain/repo"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/infrastructure/db/dbgen"
	"dreamkast-weaver/internal/pkg/stacktrace"
)

type TrailMapStampsRepoImpl struct {
	q *dbgen.Queries
}

func NewTrailMapStampsRepo(db dbgen.DBTX) repo.TrailMapStampsRepo {
	q := dbgen.New(db)
	return &TrailMapStampsRepoImpl{q}
}

var _ repo.TrailMapStampsRepo = (*TrailMapStampsRepoImpl)(nil)

func (r *TrailMapStampsRepoImpl) Get(ctx context.Context, confName value.ConfName, profileID value.ProfileID) (*dmodel.StampChallenges, error) {
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

func (r *TrailMapStampsRepoImpl) Upsert(ctx context.Context, confName value.ConfName, profileID value.ProfileID, scs *dmodel.StampChallenges) error {
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
