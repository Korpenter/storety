// Package postgres implements the database operations for the postgres database.
package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// conn is an interface that wraps the methods of pgxpool.Pool for database operations.
type conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Ping(ctx context.Context) error
	Close()
}

// DB is a wrapper around a pgxpool.Pool that implements conn interface for database operations.
type DB struct {
	conn
}

// NewDB creates a new DB instance with the given connection string.
func NewDB(connString string) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return &DB{conn: pool}, nil
}

// Ping checks if the connection to the database is still alive.
func (d *DB) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}

// commitTx commits or rolls back a transaction depending on the error.
// If there is an error, it rolls back the transaction, otherwise, it commits the transaction.
func (d *DB) commitTx(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
