// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package repo

import (
	"context"
	"database/sql"
	"time"
)

const insertCfpVote = `-- name: InsertCfpVote :exec
INSERT INTO cfp_votes (
  conference_name,
  talk_id,
  client_ip,
  created_at
) VALUES ( 
  ?, ?, ?, now()
)
`

type InsertCfpVoteParams struct {
	ConferenceName string
	TalkID         int32
	ClientIp       sql.NullString
}

func (q *Queries) InsertCfpVote(ctx context.Context, arg InsertCfpVoteParams) error {
	_, err := q.db.ExecContext(ctx, insertCfpVote, arg.ConferenceName, arg.TalkID, arg.ClientIp)
	return err
}

const listCfpVotes = `-- name: ListCfpVotes :many
SELECT conference_name, talk_id, created_at, client_ip FROM cfp_votes
WHERE
  conference_name = ? AND
  created_at > ? AND
  created_at < ?
`

type ListCfpVotesParams struct {
	ConferenceName string
	Start          time.Time
	End            time.Time
}

func (q *Queries) ListCfpVotes(ctx context.Context, arg ListCfpVotesParams) ([]CfpVote, error) {
	rows, err := q.db.QueryContext(ctx, listCfpVotes, arg.ConferenceName, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CfpVote
	for rows.Next() {
		var i CfpVote
		if err := rows.Scan(
			&i.ConferenceName,
			&i.TalkID,
			&i.CreatedAt,
			&i.ClientIp,
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
