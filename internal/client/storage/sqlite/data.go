package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/constants"
	"time"
)

// CreateData creates a new data entry in the database.
func (d *DB) CreateData(ctx context.Context, data *models.Data) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer d.commitTx(tx, err)
	res, err := d.conn.ExecContext(ctx, createData, data.ID, data.Name, data.Type, data.Content, time.Now())
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
	res, err := d.conn.ExecContext(ctx, deleteDataByName, time.Now(), name)
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

// GetSyncData retrieves id and last updated timestamp for all entries that need syncing.
func (d *DB) GetSyncData(ctx context.Context) ([]models.Data, []models.Data, time.Time, error) {
	var lastSynced time.Time
	var newData []models.Data
	var deleteData []models.Data
	var name, dataType sql.NullString
	err := d.conn.QueryRowContext(ctx, getSyncTimstamp).Scan(&lastSynced)
	rows, err := d.conn.QueryContext(ctx, getSyncData, lastSynced)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, time.Time{}, constants.ErrNoData
		}
		return nil, nil, time.Time{}, err
	}
	for rows.Next() {
		var data models.Data
		err = rows.Scan(&data.ID, &name, &dataType, &data.Content, &data.UpdatedAt, &data.Deleted)
		if err != nil {
			return nil, nil, time.Time{}, err
		}
		data.Name = name.String
		data.Type = dataType.String
		if data.Deleted {
			deleteData = append(deleteData, data)
			continue
		}
		newData = append(newData, data)
	}
	return newData, deleteData, lastSynced, nil
}

func (d *DB) UpdateSyncData(ctx context.Context, syncedNewData []models.Data, updatedData []models.Data) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer d.commitTx(tx, err)
	for _, v := range syncedNewData {
		_, err = d.conn.ExecContext(ctx, setSyncedStatus, v.ID)
		if err != nil {
			return err
		}
	}

	for _, v := range updatedData {
		_, err = d.conn.ExecContext(ctx, updateData, v.ID, v.Name, v.Type, v.Content, v.UpdatedAt, v.Deleted)
		if err != nil {
			return err
		}
	}
	_, err = d.conn.ExecContext(ctx, setLastSyncedTime, time.Now())
	if err != nil {
		return err
	}
	return nil
}
