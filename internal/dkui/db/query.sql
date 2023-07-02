-- name: ListViewEvents :many
SELECT
  *
FROM
  view_events
WHERE
  conference_name = ?
  AND profile_id = ?;

-- name: InsertViewEvents :exec
INSERT INTO
  view_events (profile_id, conference_name, track_id, talk_id, slot_id, viewing_seconds, created_at)
VALUES
  (?, ?, ?, ?, ?, ?, NOW());

-- name: GetTrailmapStamps :one
SELECT
  *
FROM
  trailmap_stamps
WHERE
  conference_name = ?
  AND profile_id = ?;

-- name: UpsertTrailmapStamp :exec
REPLACE
  trailmap_stamps (conference_name, profile_id, stamps)
VALUES
  (?, ?, ?);

-- name: UpsertViewerCount :exec
REPLACE
  viewer_counts (conference_name, track_id, channel_arn, track_name, count, updated_at)
VALUES
  (?, ?, ?, ?, ?, NOW());

-- name: ListViewerCount :many
SELECT
  *
FROM
  viewer_counts
WHERE
  conference_name = ?;
