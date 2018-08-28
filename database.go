package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// database connection used for writes
var DB_CONN *sql.DB

// TABLES_SQL defines the main database tables
// and trigger functions.
var TABLES_SQL = `
    CREATE TABLE IF NOT EXISTS assets (
        asset_id TEXT NOT NULL PRIMARY KEY,
        file_name TEXT,
        create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        update_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_deleted BOOLEAN DEFAULT f
    );

    CREATE INDEX IF NOT EXISTS assets__idx ON assets(asset_id);

    CREATE TRIGGER IF NOT EXISTS assets__update
        AFTER UPDATE ON assets FOR EACH ROW
    BEGIN
        UPDATE assets SET update_at=CURRENT_TIMESTAMP
            WHERE timestamp=OLD.timestamp;
    END;
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

func Insert(file_name, asset_id string) error {
	tx, err := DB_CONN.Begin()
	if nil != err {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR REPLACE INTO assets(
            file_name,
            asset_id
        )
        VALUES (?, ?)`)
	if nil != err {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(file_name, asset_id)
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

func Select(asset_id string) (string, error) {
	conn, _ := OpenDb()
	defer conn.Close()

	query := fmt.Sprintf(`
        SELECT
            %v
        FROM
            assets
        WHERE
            asset_id = ?
    `, AssetSQL)

	row := conn.QueryRow(query, asset_id)

	var result string
	err := row.Scan(&result)
	return result, err
}

func OpenDb() (*sql.DB, error) {
	db_conn, err := sql.Open("sqlite3", "assets.db?cache=shared&mode=rwc&_busy_timeout=50000000")
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	return db_conn, err
}

func init() {
	DB_CONN, _ = OpenDb()
	makeTables()
}
