// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package repo

import (
	"database/sql"
	"time"
)

type CfpVote struct {
	ConferenceName string
	TalkID         int32
	CreatedAt      time.Time
	ClientIp       sql.NullString
}

type SchemaMigration struct {
	Version string
}
