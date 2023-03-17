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
	err = tx.QueryRow(ctx, createUser, user.ID, user.Login, user.Password, user.Salt).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrUserExists
		}
		return err
	}
	return nil
}

// GetUserDataByName retrieves the user id, password and salt for a specific user.
// It returns the user's UUID, password, and any error that occurs.
func (d *DB) GetUserDataByName(ctx context.Context, username string) (uuid.UUID, string, string, error) {
	var password, salt string
	var id uuid.UUID
	err := d.conn.QueryRow(ctx, getUserDataByName, username).Scan(&id, &password, &salt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, "", "", constants.ErrUserNotFound
		}
		return uuid.Nil, "", "", err
	}
	return id, password, "", nil
}
