package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Ping(ctx context.Context) error
	Close()
}

type DB struct {
	conn
}

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

func (d *DB) Close() {
	d.conn.Close()
}

func (d *DB) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}

func (d *DB) commitTx(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
