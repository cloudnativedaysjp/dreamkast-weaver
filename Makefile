createmigration:
	dbmate new table

migrate:
	dbmate up

migratedown:
	dbmate down

.PHONY: migrate migratedown migrate