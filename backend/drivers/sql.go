package drivers

import (
	"github.com/inconshreveable/log15"
	"github.com/juju/errors"
	"github.com/samber/do"
	"github.com/xo/dburl"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/tooling/di"
	"github.com/khwong-c/wtcode/tooling/log"
)

type SQLTarget string

const (
	SQLTargetDefault SQLTarget = "default"
)

type SQLType string

const (
	SQLTypePostgres SQLType = "postgres"
	SQLTypeSQLite   SQLType = "sqlite"
)

type SQL interface {
	DB() *gorm.DB
}

type sql struct {
	DBType SQLType

	injector *do.Injector
	db       *gorm.DB
	logger   log15.Logger
}

func (s *sql) Shutdown() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return errors.Trace(err)
	}
	return sqlDB.Close()
}

func (s *sql) DB() *gorm.DB {
	return s.db.Session(&gorm.Session{PrepareStmt: true})
}

func getConnURL(injector *do.Injector, target string) (string, error) {
	cfg := di.InvokeOrProvide(injector, config.LoadConfig)
	switch target {
	case string(SQLTargetDefault):
		return cfg.SQLTarget.Default, nil
	default:
		return "", errors.NotFoundf("SQL target not found: %s", target)
	}
}

func createDialector(connStr string) (gorm.Dialector, SQLType, error) {
	url, err := dburl.Parse(connStr)
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	dsn := url.DSN
	switch url.Driver {
	case "postgres":
		return postgres.Open(dsn), SQLTypePostgres, nil
	case "sqlite3":
		return sqlite.Open(dsn), SQLTypeSQLite, nil
	default:
		return nil, "", errors.NotImplemented
	}
}

func DialSQL(injector *do.Injector, target string) (SQL, error) {
	logger := log.NewLogger("sql").New("target", target)
	connStr, err := getConnURL(injector, target)
	if err != nil {
		return nil, errors.Trace(err)
	}

	dialector, dbType, err := createDialector(connStr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &sql{
		DBType:   dbType,
		injector: injector,
		db:       db,
		logger:   logger,
	}, nil
}
