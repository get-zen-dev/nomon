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
	Sql  *sql.DB
	Stmt *sql.Stmt
}

// NewDB creates new database
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
		Sql:  sqlDB,
		Stmt: stmt,
	}
	return &db, nil
}

// Add adds row to the database
func (db *DB) Add(stat ServerStatus) error {
	tx, err := db.Sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.Stmt).Exec(stat.Time, stat.CPUStatus, stat.RAMStatus, stat.DiskStatus)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Close closes the database and statement
func (db *DB) Close() error {
	defer func() {
		db.Stmt.Close()
		db.Sql.Close()
	}()

	return nil
}

// PrintValues prints all rows from database
func (db *DB) PrintValues() {
	rows, err := db.Sql.Query("SELECT * FROM serverStatus")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	stats := []ServerStatus{}
	for rows.Next() {
		stat := ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUStatus, &stat.RAMStatus, &stat.DiskStatus)
		if err != nil {
			fmt.Println(err)
			continue
		}
		stats = append(stats, stat)
	}
	for _, stat := range stats {
		fmt.Println(stat.Time, stat.CPUStatus, stat.RAMStatus, stat.DiskStatus)
	}
}
