package sqlhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	_ "github.com/go-sql-driver/mysql"
)

const (
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBEndpoint = "DB_ENDPOINT"
	envDBPort     = "DB_PORT"
)

// SqlHelper provides sql helper methods like transaction and batch.
type SqlHelper struct {
	db *sql.DB
}

// SqlOption is a parameter set for SqlHelper.
type SqlOption struct {
	User     string
	Password string
	Endpoint string
	Port     string
	DbName   string
}

func NewOptionFromEnv() *SqlOption {
	return &SqlOption{
		User:     os.Getenv(envDBUser),
		Password: os.Getenv(envDBPassword),
		Endpoint: os.Getenv(envDBEndpoint),
		Port:     os.Getenv(envDBPort),
		DbName:   "dkui",
	}
}

func (o *SqlOption) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.User, validation.Required),
		validation.Field(&o.Password, validation.Required),
		validation.Field(&o.Endpoint, validation.Required),
		validation.Field(&o.Port, validation.Required),
		validation.Field(&o.DbName, validation.Required),
	)
}

// NewSqlHelper creates and sets up Database connection and returns it.
func NewSqlHelper(opt *SqlOption) (*SqlHelper, error) {
	// defaulting to env vars
	if err := opt.Validate(); err != nil {
		return nil, fmt.Errorf("invalid sql option: %w", err)
	}

	info := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?parseTime=true&loc=UTC",
		opt.User,
		opt.Password,
		opt.Endpoint,
		opt.Port,
		opt.DbName,
	)
	db, err := sql.Open("mysql", info)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &SqlHelper{db: db}, nil
}

// DB provides sql.DB.
func (s *SqlHelper) DB() *sql.DB {
	return s.db
}

// RunTX runs transaction and provides sql.Tx to the given callback.
// It performs commit only when no errors have occurred in the callback.
// Rollback will be performed any time but it effects at the failure case only.
func (s *SqlHelper) RunTX(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if e := tx.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			err = errors.Join(err, e)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return fmt.Errorf("in transaction: %w", err)
	}

	err = tx.Commit()
	return err
}

type _keySqlHelper struct{}

func WithSqlHelper(parent context.Context, h *SqlHelper) context.Context {
	return context.WithValue(parent, _keySqlHelper{}, h)
}

func SqlHelperFromContext(ctx context.Context) *SqlHelper {
	v, ok := ctx.Value(_keySqlHelper{}).(*SqlHelper)
	if !ok {
		return nil
	} else {
		return v
	}
}

func testSqlOpt(dbName string) *SqlOption {
	return &SqlOption{
		User:     "user",
		Password: "password",
		Endpoint: "127.0.0.1",
		Port:     "13306",
		DbName:   dbName,
	}
}

func NewTestSqlHelper(dbName string) *SqlHelper {
	sqlh, err := NewSqlHelper(testSqlOpt(dbName))
	if err != nil {
		log.Fatal(err)
	}
	return sqlh
}
