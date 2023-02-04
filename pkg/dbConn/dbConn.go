package dbConn

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

type ServerStatus struct {
	Time       time.Time
	CPUStatus  float64
	RAMStatus  float64
	DiskStatus float64
}

type DB struct {
	sql    *sql.DB
	stmt   *sql.Stmt
	buffer []serverStatus
	mutex  *sync.Mutex
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
		INSERT INTO trades (
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
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]serverStatus, 0, 32),
	}
	return &db, nil
}

func (db *DB) Add(stat serverStatus) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("server status buffer is full")
	}

	db.buffer = append(db.buffer, stat)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush buffer: %w", err)
		}
	}

	return nil
}

func (db *DB) Flush() error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	for _, stat := range db.buffer {
		_, err := tx.Stmt(db.stmt).Exec(stat.Time, stat.CPUStatus, stat.RAMStatus, stat.DiskStatus)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *DB) Close() error {
	defer func() {
		db.stmt.Close()
		db.sql.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
