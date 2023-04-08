-- name: ListCfpVotes :many
SELECT * FROM cfp_votes;

-- name: InsertCfpVote :exec
INSERT INTO cfp_votes (
  conference_name,
  talk_id,
  dt
) VALUES ( 
  ?, ?, ?
);

-- name: ListCfpVoteByConferenceName :many
SELECT * FROM cfp_votes
WHERE conference_name = ?;

