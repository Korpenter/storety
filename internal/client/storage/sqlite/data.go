package sqlite

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/google/uuid"
	"strings"
	"time"
)

// CreateData creates a new data entry in the database.
func (d *DB) CreateData(ctx context.Context, data *models.Data) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer d.commitTx(tx, err)
	res, err := tx.ExecContext(ctx, createData, data.ID, data.Name, data.Type, data.Content, time.Now().UTC())
	if affected, _ := res.RowsAffected(); affected == 0 || err != nil {
		return errors.Join(constants.ErrCreateData, err)
	}
	return nil
}

// GetDataContentByName retrieves the content and content type of data by name for a specific user.
func (d *DB) GetDataContentByName(ctx context.Context, name string) ([]byte, string, error) {
	var content []byte
	var contentType string
	err := d.conn.QueryRowContext(ctx, getDataContentByName, name).Scan(&content, &contentType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", constants.ErrGetData
		}
		return nil, "", err
	}
	return content, contentType, nil
}

// DeleteDataByName deletes a data entry by name for a specific user.
func (d *DB) DeleteDataByName(ctx context.Context, name string) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer d.commitTx(tx, err)
	res, err := tx.ExecContext(ctx, deleteDataByName, time.Now().UTC(), name)
	if affected, _ := res.RowsAffected(); affected == 0 || err != nil {
		return errors.Join(constants.ErrDeleteData, err)
	}
	return nil
}

// GetAllDataInfo retrieves all data info (name, type) for a specific user.
func (d *DB) GetAllDataInfo(ctx context.Context) ([]models.DataInfo, error) {
	var list []models.DataInfo
	rows, err := d.conn.QueryContext(ctx, getAllDataInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	defer rows.Close()
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

// GetNewData retrieves id and last updated timestamp for all entries that were never synced.
func (d *DB) GetNewData(ctx context.Context) ([]models.Data, error) {
	var newData []models.Data
	rows, err := d.conn.QueryContext(ctx, getNewData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data models.Data
		var name, dataType sql.NullString
		err = rows.Scan(&data.ID, &name, &dataType, &data.Content, &data.UpdatedAt, &data.Deleted)
		if err != nil {
			return nil, err
		}
		data.Name = name.String
		data.Type = dataType.String
		newData = append(newData, data)
	}
	return newData, nil
}

func (d *DB) SetSyncedStatus(ctx context.Context, newData []models.Data) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer d.commitTx(tx, err)
	for _, v := range newData {
		_, err = tx.ExecContext(ctx, setSyncedStatus, v.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) SyncBatch(ctx context.Context, syncBatch []models.Data) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer d.commitTx(tx, err)
	for _, v := range syncBatch {
		_, err = tx.ExecContext(ctx, insertOrReplaceData, v.ID, v.Name, v.Type, v.Content, v.UpdatedAt, v.Deleted)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) GetBatch(ctx context.Context, ids []uuid.UUID) ([]models.Data, error) {
	var batch []models.Data
	query := getBatch + strings.Repeat(", ?", len(ids)-1) + `)`
	intIds := make([]interface{}, len(ids))
	for i, v := range ids {
		intIds[i] = v
	}
	rows, err := d.conn.QueryContext(ctx, query, intIds...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data models.Data
		var name, dataType sql.NullString
		err = rows.Scan(&data.ID, &name, &dataType, &data.Content, &data.UpdatedAt, &data.Deleted)
		if err != nil {
			return nil, err
		}
		data.Name = name.String
		data.Type = dataType.String
		batch = append(batch, data)
	}
	return batch, nil
}

// GetSyncData retrieves id and last updated timestamp and content hash for all entries that were ever synced.
func (d *DB) GetSyncData(ctx context.Context) ([]models.SyncData, error) {
	var syncData []models.SyncData
	rows, err := d.conn.QueryContext(ctx, getSyncData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		hasher := md5.New()
		var data models.SyncData
		var content []byte
		err = rows.Scan(&data.ID, &content, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hasher.Write(content)
		data.Hash = hex.EncodeToString(hasher.Sum(nil))
		syncData = append(syncData, data)
	}
	return syncData, nil
}
