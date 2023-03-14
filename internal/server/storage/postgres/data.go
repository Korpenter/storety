package postgres

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (d *DB) CreateData(ctx context.Context, data *models.Data) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := d.conn.Exec(ctx, createData, data.ID, data.UserID, data.Name, data.Type, data.Content)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(storage.ErrCreatingData, err)
	}
	res, err = d.conn.Exec(ctx, updateDataVersion, data.UserID)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(storage.ErrUpdatingVersion, err)
	}
	return nil
}

func (d *DB) GetDataContentByName(ctx context.Context, userID uuid.UUID, name string) ([]byte, error) {
	var content []byte
	err := d.conn.QueryRow(ctx, getDataContentByName, userID, name).Scan(&content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrGettingData
		}
		return nil, err
	}
	return content, nil
}

func (d *DB) DeleteDataByName(ctx context.Context, userID uuid.UUID, name string) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := d.conn.Exec(ctx, deleteDataByName, userID, name)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(storage.ErrDeletingData, err)
	}
	res, err = d.conn.Exec(ctx, updateDataVersion, userID)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(storage.ErrCreatingSession, err)
	}
	return nil
}

func (d *DB) GetAllDataInfo(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error) {
	var list []models.DataInfo
	rows, err := d.conn.Query(ctx, getAllDataInfo, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNoData
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
