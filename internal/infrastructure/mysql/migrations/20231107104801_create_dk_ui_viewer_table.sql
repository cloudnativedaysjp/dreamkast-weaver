-- migrate:up
CREATE TABLE
  track_viewer (
    created_at DATETIME(3) NOT NULL,
    track_name CHAR(1) NOT NULL,
    profile_id INT NOT NULL,
    PRIMARY KEY (created_at)
  );

-- migrate:down
DROP TABLE track_viewer
