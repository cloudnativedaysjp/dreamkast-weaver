-- migrate:up
ALTER TABLE track_viewer ADD COLUMN talk_id INT NOT NULL;

-- migrate:down
ALTER TABLE track_viewer DROP COLUMN talk_id
