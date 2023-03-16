package postgres

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateSession creates a new session entry in the database.
func (d *DB) CreateSession(ctx context.Context, session *models.Session, oldSession *models.Session) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := d.conn.Exec(ctx, createNewSession, session.ID, session.UserID, session.AuthToken, session.RefreshToken)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrCreateSession, err)
	}
	if oldSession != nil {
		res, err = d.conn.Exec(ctx, deleteOldSession, oldSession.ID, oldSession.RefreshToken)
		if res.RowsAffected() == 0 || err != nil {
			return errors.Join(constants.ErrDeleteSession, err)
		}
	}
	return nil
}

// GetSession retrieves the user id for a specific session.
func (d *DB) GetSession(ctx context.Context, sessionID uuid.UUID, refreshToken string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := d.conn.QueryRow(ctx, getUserBySession, sessionID, refreshToken).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, constants.ErrSessionNotFound
		}
		return uuid.Nil, err
	}
	return userID, nil
}
