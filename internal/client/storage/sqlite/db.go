// Package sqlite implements the database operations for the SQLite database.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

// conn is an interface that wraps the methods of sql.DB for database operations.
type conn interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Close() error
}

// DB is a wrapper around a sql.DB that implements conn interface for database operations.
type DB struct {
	conn
}

// NewDB creates a new DB instance with the given SQLite database file path.
func NewDB(databasePath, username string) (*DB, error) {
	database := filepath.Join(databasePath, username+".db")
	_, err := os.Stat(database)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(database)
			if err != nil {
				return nil, err
			}
			file.Close()
		} else {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", database))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createTableData)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	return &DB{conn: db}, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}

// commitTx commits or rolls back a transaction depending on the error.
// If there is an error, it rolls back the transaction, otherwise, it commits the transaction.
func (d *DB) commitTx(tx *sql.Tx, err error) {
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
