package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB_CONN *sql.DB

// TABLES_SQL defines the main database tables
// and trigger functions.
var TABLES_SQL = `
    CREATE TABLE IF NOT EXISTS uploads (
        uid TEXT NOT NULL PRIMARY KEY,
        file_name TEXT,
        create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        update_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS uploads__idx ON uploads(uid);
    CREATE TRIGGER IF NOT EXISTS uploads__update
        AFTER
        UPDATE
        ON uploads
        FOR EACH ROW
    BEGIN
        UPDATE uploads SET update_at=CURRENT_TIMESTAMP WHERE timestamp=OLD.timestamp;
    END;
`

var ASSET_SQL = `
	'{'||
		'"uid": ' ||  uid ||','||
		'"file_name": "' ||  file_name ||'",'||
		'"create_at": "' ||  create_at ||'",'||
		'"update_at": "' ||  update_at
	|| '}'
`

// MakeTables creates database tables and triggers.
func makeTables() (err error) {
	logger.Debug("create database tables")
	_, err = DB_CONN.Exec(TABLES_SQL)
	if err != nil {
		logger.Error(err)
		panic("failed to create database tables")
		return
	}
	return
}

func Insert(file_name, uid string) error {
	tx, err := DB_CONN.Begin()
	if nil != err {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR REPLACE INTO uploads(
            file_name,
            uid
        )
        VALUES (?, ?)`)
	if nil != err {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(file_name, uid)
	if nil != err {
		return err
	}

	err_commit := tx.Commit()
	if nil != err {
		err_rollback := tx.Rollback()
		if nil != err_rollback {
			return err_rollback
		}
		return err_commit
	}

	return nil
}

func init() {
	var err error
	DB_CONN, err = sql.Open("sqlite3", "assets.db?cache=shared&mode=rwc&_busy_timeout=50000000")
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	makeTables()
}
