-- migrate:up
ALTER TABLE track_viewer DROP PRIMARY KEY, ADD PRIMARY KEY(created_at,profile_id);

-- migrate:down
ALTER TABLE track_viewer DROP PRIMARY KEY, ADD PRIMARY KEY(created_at);
