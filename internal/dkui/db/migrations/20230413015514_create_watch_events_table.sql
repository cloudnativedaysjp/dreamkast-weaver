-- migrate:up
CREATE TABLE
  watch_events (
    conference_name VARCHAR(32) NOT NULL,
    profile_id INT NOT NULL,
    track_id INT NOT NULL,
    talk_id INT NOT NULL,
    slot_id INT NOT NULL,
    viewing_seconds INT NOT NULL,
    created_at DATETIME NOT NULL
  );

-- migrate:down
DROP TABLE watch_events