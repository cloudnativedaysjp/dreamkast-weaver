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

const (
	dbName         = "/dkui"
	dbMigrationDir = "internal/infrastructure/db/migrations"
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
		mdp := filepath.Join(wd, dbMigrationDir)
		if _, err := os.Stat(mdp); err != nil {
			log.Fatal(err)
		}

		u := &url.URL{
			Scheme: "mysql",
			User:   url.UserPassword(user, pw),
			Host:   fmt.Sprintf("%s:%s", endpoint, port),
			Path:   dbName,
		}

		dbm := dbmate.New(u)
		dbm.MigrationsDir = []string{filepath.Join(wd, dbMigrationDir)}
		if err := dbm.CreateAndMigrate(); err != nil {
			log.Fatal(err)
		}

		log.Printf("CreateAndMigrate for %s is complete", dbName)
	},
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
