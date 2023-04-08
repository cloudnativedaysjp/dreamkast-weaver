createmigration:
	dbmate new table

migrate:
	dbmate up

migratedown:
	dbmate down

gplgen-generate:
	go run github.com/99designs/gqlgen generate

.PHONY: migrate migratedown migrate