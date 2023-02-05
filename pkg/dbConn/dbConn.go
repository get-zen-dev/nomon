package dbConn

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ServerStatus struct {
	Time       time.Time
	CPUStatus  float64
	RAMStatus  float64
	DiskStatus float64
}

type DB struct {
	sql  *sql.DB
	stmt *sql.Stmt
}

func NewDB(dbFile string) (*DB, error) {

	schemaSQL := `
		CREATE TABLE IF NOT EXISTS serverStatus (
			time TIMESTAMP,
			cpustatus FLOAT,
			ramstatus FLOAT,
			diskstatus FLOAT
		);

		CREATE INDEX IF NOT EXISTS status_time ON serverStatus(time);
		`
	insertSQL := `
		INSERT INTO serverStatus (
			time, cpustatus, ramstatus, diskstatus
		) VALUES (
			?, ?, ?, ?
		)`

	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		return nil, err
	}

	stmt, err := sqlDB.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}
	db := DB{
		sql:  sqlDB,
		stmt: stmt,
	}
	return &db, nil
}

func (db *DB) Add(stat ServerStatus) error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.stmt).Exec(stat.Time, stat.CPUStatus, stat.RAMStatus, stat.DiskStatus)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) Close() error {
	defer func() {
		db.stmt.Close()
		db.sql.Close()
	}()

	return nil
}

func (db *DB) PrintValues() {
	rows, err := db.sql.Query("SELECT time, cpustatus, ramstatus, diskstatus FROM serverStatus")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println(rows.Scan())
	}
}
