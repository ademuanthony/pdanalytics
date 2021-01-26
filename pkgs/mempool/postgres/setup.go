package postgres

//go:generate sqlboiler --wipe psql --no-hooks --no-auto-timestamps

import (
	"context"
	"fmt"
)

var (
	createMempoolTable = `CREATE TABLE IF NOT EXISTS mempool (
		time timestamp,
		first_seen_time timestamp,
		number_of_transactions INT,
		voters INT,
		tickets INT,
		revocations INT,
		size INT,
		total_fee FLOAT8,
		total FLOAT8,
		PRIMARY KEY (time)
	);`

	createMempoolDayBinTable = `CREATE TABLE IF NOT EXISTS mempool_bin (
		time INT8,
		bin VARCHAR(25),
		number_of_transactions INT,
		size INT,
		total_fee FLOAT8,
		PRIMARY KEY (time,bin)
	);`

	lastMempoolBlockHeight = `SELECT last_block_height FROM mempool ORDER BY last_block_height DESC LIMIT 1`
	lastMempoolEntryTime   = `SELECT time FROM mempool ORDER BY time DESC LIMIT 1`
)

func (db *PgDb) CreateTables(ctx context.Context) error {
	if !db.mempoolDataTableExits() {
		if err := db.createMempoolDataTable(); err != nil {
			return err
		}
	}
	if !db.mempoolBinDataTableExits() {
		if err := db.createMempoolDayBinTable(); err != nil {
			return err
		}
	}
}

func (pg *PgDb) createMempoolDataTable() error {
	_, err := pg.db.Exec(createMempoolTable)
	return err
}

func (pg *PgDb) createMempoolDayBinTable() error {
	_, err := pg.db.Exec(createMempoolDayBinTable)
	return err
}

func (pg *PgDb) mempoolDataTableExits() bool {
	exists, _ := pg.tableExists("mempool")
	return exists
}

func (pg *PgDb) mempoolBinDataTableExits() bool {
	exists, _ := pg.tableExists("mempool_bin")
	return exists
}

func (pg *PgDb) tableExists(name string) (bool, error) {
	rows, err := pg.db.Query(`SELECT relname FROM pg_class WHERE relname = $1`, name)
	if err == nil {
		defer func() {
			if e := rows.Close(); e != nil {
				log.Error("Close of Query failed: ", e)
			}
		}()
		return rows.Next(), nil
	}
	return false, err
}

func (pg *PgDb) DropTables() error {

	// mempool
	if err := pg.dropTable("mempool"); err != nil {
		return err
	}

	// mempool_bin
	if err := pg.dropTable("mempool_bin"); err != nil {
		return err
	}

	return nil
}

func (pg *PgDb) ClearCache() error {
	// mempool_bin
	if err := pg.dropTable("mempool_bin"); err != nil {
		return err
	}

	return nil
}

func (pg *PgDb) dropTable(name string) error {
	log.Tracef("Dropping table %s", name)
	_, err := pg.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name))
	return err
}