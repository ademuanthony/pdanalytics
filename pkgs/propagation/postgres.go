package propagation

import (
	"database/sql"
	"time"

	"github.com/planetdecred/pdanalytics/dbhelpers"
	"github.com/volatiletech/sqlboiler/boil"
)

type PgDb struct {
	db           *sql.DB
	queryTimeout time.Duration
}

type logWriter struct{}

func (l logWriter) Write(p []byte) (n int, err error) {
	log.Debug(string(p))
	return len(p), nil
}

func NewPgDb(host, port, user, pass, dbname string, debug bool) (*PgDb, error) {
	db, err := dbhelpers.Connect(host, port, user, pass, dbname)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	if debug {
		boil.DebugMode = true
		boil.DebugWriter = logWriter{}
	}
	return &PgDb{
		db:           db,
		queryTimeout: time.Second * 30,
	}, nil
}

func (pg *PgDb) Close() error {
	log.Trace("Closing postgresql connection")
	return pg.db.Close()
}
