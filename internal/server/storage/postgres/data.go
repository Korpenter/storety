package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateData creates a new data entry in the database.
func (d *DB) CreateData(ctx context.Context, userID uuid.UUID, data *models.Data) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := d.conn.Exec(ctx, createData, data.ID, userID, data.Name, data.Type, data.Content, data.UpdatedAt, data.Deleted)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrCreateData, err)
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
			return nil, "", constants.ErrGetData
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
		return errors.Join(constants.ErrDeleteData, err)
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

// SyncData retrieves all data entries that have been updated since the last sync for a specific user
// and applies updates from client.
func (d *DB) SyncData(ctx context.Context, userID uuid.UUID, syncData models.SyncData) ([]models.Data, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer d.commitTx(ctx, tx, err)
	if len(syncData.CreateData) > 0 {
		batch := &pgx.Batch{}
		for _, data := range syncData.CreateData {
			batch.Queue(createData, data.ID, userID, data.Name, data.Type, data.Content, data.UpdatedAt, data.Deleted)
		}
		br := d.conn.SendBatch(ctx, batch)
		if res, err := br.Exec(); res.RowsAffected() == 0 || err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.Join(constants.ErrCreateData, err)
			}
			return nil, err
		}
	}
	if len(syncData.DeleteData) > 0 {
		batch := &pgx.Batch{}
		for _, data := range syncData.DeleteData {
			batch.Queue(deleteDataByID, data.ID, userID, data.UpdatedAt)
		}
		br := d.conn.SendBatch(ctx, batch)
		if res, err := br.Exec(); res.RowsAffected() == 0 || err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.Join(constants.ErrCreateData, err)
			}
			return nil, err
		}
	}
	var updates []models.Data
	rows, err := d.conn.Query(ctx, getDataBySyncTime, userID, syncData.LastSync)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	for rows.Next() {
		var data models.Data
		var name, dataType sql.NullString
		err = rows.Scan(&data.ID, &name, &dataType, &data.Content, &data.UpdatedAt, &data.Deleted)
		if err != nil {
			return nil, err
		}
		data.Name = name.String
		data.Type = dataType.String
		updates = append(updates, data)
	}
	return updates, nil
}
