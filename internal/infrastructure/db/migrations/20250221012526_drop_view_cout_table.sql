-- migrate:up
DROP TABLE viewer_counts

-- migrate:down
CREATE TABLE
  viewer_counts (
    conference_name VARCHAR(32) NOT NULL,
    track_id INT NOT NULL,
    channel_arn VARCHAR(128) NOT NULL,
    track_name VARCHAR(32) NOT NULL,
    count BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (track_id)
  );
