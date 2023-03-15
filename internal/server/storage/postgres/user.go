package postgres

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateUser creates a new user entry in the database.
func (d *DB) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	err = tx.QueryRow(ctx, createUser, user.ID, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrUserExists
		}
		return err
	}
	return nil
}

// GetIdPassByName retrieves the user id and password for a specific user.
func (d *DB) GetIdPassByName(ctx context.Context, username string) (uuid.UUID, string, error) {
	var password string
	var id uuid.UUID
	err := d.conn.QueryRow(ctx, getUserByName, username).Scan(&id, &password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, "", constants.ErrUserNotFound
		}
		return uuid.Nil, "", err
	}
	return id, password, nil
}
