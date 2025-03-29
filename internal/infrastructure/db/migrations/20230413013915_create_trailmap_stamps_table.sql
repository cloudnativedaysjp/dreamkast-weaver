-- migrate:up
CREATE TABLE
  trailmap_stamps (
    conference_name VARCHAR(32) NOT NULL,
    profile_id INT NOT NULL,
    stamps JSON NOT NULL,
    PRIMARY KEY (conference_name, profile_id)
  );

-- migrate:down
DROP TABLE trailmap_stamps