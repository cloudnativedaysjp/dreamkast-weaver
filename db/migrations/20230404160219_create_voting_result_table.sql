-- migrate:up
CREATE TABLE cfp_votes (
  conference_name VARCHAR(36) NOT NULL,
  talk_id INT NOT NULL,
  dt DATETIME NOT NULL
);

-- migrate:down
DROP TABLE cfp_votes