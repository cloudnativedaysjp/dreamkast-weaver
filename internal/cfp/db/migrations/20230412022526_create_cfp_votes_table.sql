-- migrate:up
CREATE TABLE
  cfp_votes (
    conference_name CHAR(36) NOT NULL,
    talk_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    client_ip CHAR(16)
  );

-- migrate:down
DROP TABLE cfp_votes
