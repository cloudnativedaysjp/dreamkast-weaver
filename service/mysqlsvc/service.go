package mysqlsvc

import (
	"context"
	"database/sql"
	"dreamkast-weaver/infrastructure/mysql"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"

	_ "github.com/go-sql-driver/mysql"
)

const dbName = "weaver"

type VotingResultItem struct {
	weaver.AutoMarshal
	ConferenceName string
	TarkId         int32
	Dt             time.Time
}

type impl struct {
	weaver.Implements[T]
	weaver.WithConfig[config]
	client *mysql.Queries
}

type T interface {
	InsertCfpVote(ctx context.Context) error
	ListCfpVotes(ctx context.Context) ([]VotingResultItem, error)
	//ListCfpVoteByConferenceName(ctx context.Context, conferenceName string) ([]mysql.CfpVote, error)
}

type config struct {
	Driver string `toml:"db_driver"` // Name of the database driver.
	URI    string `toml:"db_uri"`    // Database server URI.
}

func (cfg *config) Validate() error {
	if cfg.Driver != "" {
		if len(cfg.URI) == 0 {
			return fmt.Errorf("DB driver specified but not location of database")
		}
	}
	return nil
}

func (s *impl) Init(ctx context.Context) error {
	//cfg := s.Config()
	cfg := config{
		Driver: "mysql",
		URI:    "user:password@tcp(localhost:13306)/",
	}

	var db *sql.DB
	var err error

	// Ensure chat database exists.
	ensureDB := func() error {
		db_admin, err := sql.Open(cfg.Driver, cfg.URI)
		if err != nil {
			return fmt.Errorf("error opening %q URI %q: %w", cfg.Driver, cfg.URI, err)
		}
		defer db_admin.Close()
		_, err = db_admin.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		return err
	}

	if err := ensureDB(); err != nil {
		return fmt.Errorf("error creating %q database %s%s: %w", cfg.Driver, cfg.URI, dbName, err)
	}

	db, err = sql.Open(cfg.Driver, cfg.URI+dbName+"?parseTime=true")
	if err != nil {
		return fmt.Errorf("error opening %q database %s%s: %w", cfg.Driver, cfg.URI, dbName, err)
	}

	s.client = mysql.New(db)
	return nil
}

func (s *impl) InsertCfpVote(ctx context.Context) error {

	p := mysql.InsertCfpVoteParams{
		ConferenceName: "hoge",
		TalkID:         1,
		Dt:             time.Now(),
	}

	return s.client.InsertCfpVote(ctx, p)
}

func (s *impl) ListCfpVotes(ctx context.Context) ([]VotingResultItem, error) {
	res, err := s.client.ListCfpVotes(ctx)

	items := make([]VotingResultItem, 0, len(res)+1)
	for _, v := range res {
		item := VotingResultItem{
			ConferenceName: v.ConferenceName,
			TarkId:         v.TalkID,
			Dt:             v.Dt,
		}
		items = append(items, item)
	}
	return items, err
}
