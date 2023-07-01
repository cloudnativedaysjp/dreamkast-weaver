// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package repo

import (
	"context"
	"encoding/json"
)

const getTrailmapStamps = `-- name: GetTrailmapStamps :one
SELECT
  conference_name, profile_id, stamps
FROM
  trailmap_stamps
WHERE
  conference_name = ?
  AND profile_id = ?
`

type GetTrailmapStampsParams struct {
	ConferenceName string
	ProfileID      int32
}

func (q *Queries) GetTrailmapStamps(ctx context.Context, arg GetTrailmapStampsParams) (TrailmapStamp, error) {
	row := q.db.QueryRowContext(ctx, getTrailmapStamps, arg.ConferenceName, arg.ProfileID)
	var i TrailmapStamp
	err := row.Scan(&i.ConferenceName, &i.ProfileID, &i.Stamps)
	return i, err
}

const getViewerCount = `-- name: GetViewerCount :one
SELECT
  conference_name, track_id, channel_arn, track_name, count, updated_at
FROM
  viewer_counts
WHERE
  conference_name = ?
  AND track_id= ?
`

type GetViewerCountParams struct {
	ConferenceName string
	TrackID        int32
}

func (q *Queries) GetViewerCount(ctx context.Context, arg GetViewerCountParams) (ViewerCount, error) {
	row := q.db.QueryRowContext(ctx, getViewerCount, arg.ConferenceName, arg.TrackID)
	var i ViewerCount
	err := row.Scan(
		&i.ConferenceName,
		&i.TrackID,
		&i.ChannelArn,
		&i.TrackName,
		&i.Count,
		&i.UpdatedAt,
	)
	return i, err
}

const insertViewEvents = `-- name: InsertViewEvents :exec
INSERT INTO
  view_events (profile_id, conference_name, track_id, talk_id, slot_id, viewing_seconds, created_at)
VALUES
  (?, ?, ?, ?, ?, ?, NOW())
`

type InsertViewEventsParams struct {
	ProfileID      int32
	ConferenceName string
	TrackID        int32
	TalkID         int32
	SlotID         int32
	ViewingSeconds int32
}

func (q *Queries) InsertViewEvents(ctx context.Context, arg InsertViewEventsParams) error {
	_, err := q.db.ExecContext(ctx, insertViewEvents,
		arg.ProfileID,
		arg.ConferenceName,
		arg.TrackID,
		arg.TalkID,
		arg.SlotID,
		arg.ViewingSeconds,
	)
	return err
}

const listViewEvents = `-- name: ListViewEvents :many
SELECT
  conference_name, profile_id, track_id, talk_id, slot_id, viewing_seconds, created_at
FROM
  view_events
WHERE
  conference_name = ?
  AND profile_id = ?
`

type ListViewEventsParams struct {
	ConferenceName string
	ProfileID      int32
}

func (q *Queries) ListViewEvents(ctx context.Context, arg ListViewEventsParams) ([]ViewEvent, error) {
	rows, err := q.db.QueryContext(ctx, listViewEvents, arg.ConferenceName, arg.ProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ViewEvent
	for rows.Next() {
		var i ViewEvent
		if err := rows.Scan(
			&i.ConferenceName,
			&i.ProfileID,
			&i.TrackID,
			&i.TalkID,
			&i.SlotID,
			&i.ViewingSeconds,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertTrailmapStamp = `-- name: UpsertTrailmapStamp :exec
REPLACE
  trailmap_stamps (conference_name, profile_id, stamps)
VALUES
  (?, ?, ?)
`

type UpsertTrailmapStampParams struct {
	ConferenceName string
	ProfileID      int32
	Stamps         json.RawMessage
}

func (q *Queries) UpsertTrailmapStamp(ctx context.Context, arg UpsertTrailmapStampParams) error {
	_, err := q.db.ExecContext(ctx, upsertTrailmapStamp, arg.ConferenceName, arg.ProfileID, arg.Stamps)
	return err
}

const upsertViewerCount = `-- name: UpsertViewerCount :exec
REPLACE
  viewer_counts (conference_name, track_id, channel_arn, track_name, count, updated_at)
VALUES
  (?, ?, ?, ?, ?, NOW())
`

type UpsertViewerCountParams struct {
	ConferenceName string
	TrackID        int32
	ChannelArn     string
	TrackName      string
	Count          int64
}

func (q *Queries) UpsertViewerCount(ctx context.Context, arg UpsertViewerCountParams) error {
	_, err := q.db.ExecContext(ctx, upsertViewerCount,
		arg.ConferenceName,
		arg.TrackID,
		arg.ChannelArn,
		arg.TrackName,
		arg.Count,
	)
	return err
}
