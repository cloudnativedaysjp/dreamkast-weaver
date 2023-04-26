-- name: ListCfpVotes :many
SELECT * FROM cfp_votes
WHERE
  conference_name = ? AND
  created_at > sqlc.arg(start) AND
  created_at < sqlc.arg(end);

-- name: InsertCfpVote :exec
INSERT INTO cfp_votes (
  conference_name,
  talk_id,
  client_ip,
  created_at
) VALUES ( 
  ?, ?, ?, now()
);

