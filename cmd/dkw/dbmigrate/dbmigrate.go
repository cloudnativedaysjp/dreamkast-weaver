package dbmigrate

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/spf13/cobra"
)

type db struct {
	name         string
	migrationDir string
}

var (
	dbs = []db{
		{
			name:         "/dkui",
			migrationDir: "internal/dkui/db/migrations",
		},
		{
			name:         "/cfp",
			migrationDir: "internal/cfp/db/migrations",
		},
	}
)

// Cmd represents the dbmigrate command.
var Cmd = &cobra.Command{
	Use:   "dbmigrate",
	Short: "Create database and Migrate to the latest version",
	Long:  "Create database and Migrate to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		user := getEnv("DB_USER", "user")
		pw := getEnv("DB_PASSWORD", "password")
		endpoint := getEnv("DB_ENDPOINT", "localhost")
		port := getEnv("DB_PORT", "13306")

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Check for existing migration dir.
		for _, db := range dbs {
			mdp := filepath.Join(wd, db.migrationDir)
			if _, err := os.Stat(mdp); err != nil {
				log.Fatal(err)
			}
		}

		for _, db := range dbs {
			u := &url.URL{
				Scheme: "mysql",
				User:   url.UserPassword(user, pw),
				Host:   fmt.Sprintf("%s:%s", endpoint, port),
				Path:   db.name,
			}

			dbm := dbmate.New(u)
			dbm.MigrationsDir = []string{filepath.Join(wd, db.migrationDir)}
			err := dbm.CreateAndMigrate()
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("CreateAndMigrate for %s is complete", db.name)
		}
	},
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
