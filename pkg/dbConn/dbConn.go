package dbConn

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type ServerStatus struct {
	Time     time.Time
	CPUUsed  float64
	RAMUsed  uint64
	DiskUsed uint64
}

type DB struct {
	Sql  *sql.DB
	Stmt *sql.Stmt
}

// NewDB creates new database
func NewDB(dbFile string) (*DB, error) {
	log.Trace("dbConn:NewDB")
	schemaSQL := `
		CREATE TABLE IF NOT EXISTS serverStatus (
			time TIMESTAMP,
			cpuused FLOAT,
			ramused INTEGER,
			diskused INTEGER
		);

		CREATE INDEX IF NOT EXISTS status_time ON serverStatus(time);
		`
	insertSQL := `
		INSERT INTO serverStatus (
			time, cpuused, ramused, diskused
		) VALUES (
			?, ?, ?, ?
		)`
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Error("Error opening dbFile: ", err)
		return nil, err
	}
	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		log.Error("Error making SQL schema: ", err)
		return nil, err
	}
	stmt, err := sqlDB.Prepare(insertSQL)
	if err != nil {
		log.Error("Error preparing insert statement: ", err)
		return nil, err
	}
	db := DB{
		Sql:  sqlDB,
		Stmt: stmt,
	}
	return &db, nil
}

// Add adds row to the database
func (db *DB) Add(stat ServerStatus) error {
	log.Trace("dbConn:Add", stat)
	tx, err := db.Sql.Begin()
	if err != nil {
		log.Error("Error starting transaction: ", err)
		return err
	}

	_, err = tx.Stmt(db.Stmt).Exec(stat.Time, stat.CPUUsed, stat.RAMUsed, stat.DiskUsed)
	if err != nil {
		log.Error("Error executing transaction: ", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Close closes the database and statement
func (db *DB) Close() error {
	log.Trace("dbConn:Close")
	defer func() {
		db.Stmt.Close()
		db.Sql.Close()
	}()

	return nil
}

// PrintValues prints all rows from database
func (db *DB) PrintValues() {
	log.Trace("dbConn:PrintValues")
	rows, err := db.Sql.Query("SELECT * FROM serverStatus")
	if err != nil {
		log.Error("Error getting rows from db: ", err)
	}
	defer rows.Close()

	stats := []ServerStatus{}
	for rows.Next() {
		stat := ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUUsed, &stat.RAMUsed, &stat.DiskUsed)
		if err != nil {
			log.Error("Error scanning rows: ", err)
			continue
		}
		stats = append(stats, stat)
	}
	for _, stat := range stats {
		log.Info(stat.Time, stat.CPUUsed, stat.RAMUsed, stat.DiskUsed)
	}
}
