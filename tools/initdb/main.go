//go:build tools
// +build tools

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	user := getEnv("DB_USER", "user")
	pw := getEnv("DB_PASSWORD", "password")
	endpoint := getEnv("DB_ENDPOINT", "localhost")
	port := getEnv("DB_PORT", "13306")

	dbs := []struct {
		name         string
		migrationDir string
	}{
		{
			name:         "/dkui",
			migrationDir: "/internal/dkui/db/migrations",
		},
		{
			name:         "/cfp",
			migrationDir: "/internal/cfp/db/migrations",
		},
	}

	for _, db := range dbs {
		u := &url.URL{
			Scheme: "mysql",
			User:   url.UserPassword(user, pw),
			Host:   fmt.Sprintf("%s:%s", endpoint, port),
			Path:   db.name,
		}

		fmt.Printf("connecting %s ... \n", u)
		dbm := dbmate.New(u)
		dbm.MigrationsDir = db.migrationDir
		err := dbm.CreateAndMigrate()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("CreateAndMigrate done")
	}
}
