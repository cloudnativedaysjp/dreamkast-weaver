// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package repo

import (
	"encoding/json"
	"time"
)

type SchemaMigration struct {
	Version string
}

type TrailmapStamp struct {
	ConferenceName string
	ProfileID      int32
	Stamps         json.RawMessage
}

type ViewEvent struct {
	ConferenceName string
	ProfileID      int32
	TrackID        int32
	TalkID         int32
	SlotID         int32
	ViewingSeconds int32
	CreatedAt      time.Time
}

type ViewerCount struct {
	ConferenceName string
	TrackID        int32
	ChannelArn     string
	TrackName      string
	Count          int64
	UpdatedAt      time.Time
}
