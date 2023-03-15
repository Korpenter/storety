package postgres

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateData creates a new data entry in the database.
func (d *DB) CreateData(ctx context.Context, data *models.Data) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := d.conn.Exec(ctx, createData, data.ID, data.UserID, data.Name, data.Type, data.Content)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrCreatingData, err)
	}
	res, err = d.conn.Exec(ctx, updateDataVersion, data.UserID)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrUpdatingVersion, err)
	}
	return nil
}

// GetDataContentByName retrieves the content and content type of data by name for a specific user.
func (d *DB) GetDataContentByName(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error) {
	var content []byte
	var contentType string
	err := d.conn.QueryRow(ctx, getDataContentByName, name, userID).Scan(&content, &contentType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", constants.ErrGettingData
		}
		return nil, "", err
	}
	return content, contentType, nil
}

// DeleteDataByName deletes a data entry by name for a specific user.
func (d *DB) DeleteDataByName(ctx context.Context, userID uuid.UUID, name string) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := d.conn.Exec(ctx, deleteDataByName, name, userID)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrDeletingData, err)
	}
	res, err = d.conn.Exec(ctx, updateDataVersion, userID)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrCreatingSession, err)
	}
	return nil
}

// GetAllDataInfo retrieves all data info (name, type) for a specific user.
func (d *DB) GetAllDataInfo(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error) {
	var list []models.DataInfo
	rows, err := d.conn.Query(ctx, getAllDataInfo, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	for rows.Next() {
		var data models.DataInfo
		err = rows.Scan(&data.Name, &data.Type)
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return list, nil
}
