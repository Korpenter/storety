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

// CreateData implements the data service interface CreateData method.
func (d *DB) CreateData(ctx context.Context, userID uuid.UUID, data *models.Data) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := tx.Exec(ctx, createData, data.ID, userID, data.Name, data.Type, data.Content, data.UpdatedAt, data.Deleted)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrCreateData, err)
	}
	return nil
}

// GetDataContentByName implements the data service interface GetDataContentByName method.
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

// DeleteDataByName implements the data service interface DeleteDataByName method.
func (d *DB) DeleteDataByName(ctx context.Context, userID uuid.UUID, name string) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	res, err := tx.Exec(ctx, deleteDataByName, name, userID)
	if res.RowsAffected() == 0 || err != nil {
		return errors.Join(constants.ErrDeleteData, err)
	}
	return nil
}

// GetAllDataInfo implements the data service interface GetAllDataInfo method.
func (d *DB) GetAllDataInfo(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error) {
	var list []models.DataInfo
	rows, err := d.conn.Query(ctx, getAllDataInfo, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

// CreateBatch implements the DataRepository interface CreateBatch method.
func (d *DB) CreateBatch(ctx context.Context, userID uuid.UUID, dataBatch []models.Data) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	batch := &pgx.Batch{}
	for _, data := range dataBatch {
		batch.Queue(createData, data.ID, userID, data.Name, data.Type, data.Content, data.UpdatedAt, data.Deleted)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()
	if res, err := br.Exec(); res.RowsAffected() == 0 || err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(constants.ErrCreateData, err)
		}
		return err
	}
	return nil
}

// UpdateBatch implements the DataRepository interface UpdateBatch method.
func (d *DB) UpdateBatch(ctx context.Context, userID uuid.UUID, dataBatch []models.Data) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer d.commitTx(ctx, tx, err)
	batch := &pgx.Batch{}
	for _, data := range dataBatch {
		batch.Queue(updateDataByID, data.ID, userID, data.Name, data.Type, data.Content, data.Deleted, data.UpdatedAt)
	}
	br := tx.SendBatch(ctx, batch)
	defer br.Close()
	for i := 0; i < len(dataBatch); i++ {

		if res, err := br.Exec(); res.RowsAffected() == 0 || err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Join(constants.ErrUpdateData, err)
			}
			return err
		}
	}
	return nil
}

// GetNewData implements the DataRepository interface GetNewData method.
func (d *DB) GetNewData(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) ([]models.Data, error) {
	rows, err := d.conn.Query(ctx, getNewData, userID, ids)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNoData
		}
		return nil, err
	}
	defer rows.Close()
	var list []models.Data
	for rows.Next() {
		var data models.Data
		err = rows.Scan(&data.ID, &data.Name, &data.Type, &data.Content, &data.UpdatedAt, &data.Deleted)
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return list, nil
}

func (d *DB) GetDataByUpdateAndHash(ctx context.Context, userID uuid.UUID, syncData []models.SyncData) ([]models.Data, []string, error) {
	earlierBatch := &pgx.Batch{}
	laterBatch := &pgx.Batch{}
	for _, data := range syncData {
		earlierBatch.Queue(getEarlierUpdate, userID, data.ID, data.Hash, data.UpdatedAt)
		laterBatch.Queue(getLaterUpdate, userID, data.ID, data.Hash, data.UpdatedAt)
	}

	ber := d.conn.SendBatch(ctx, earlierBatch)
	defer ber.Close()
	blr := d.conn.SendBatch(ctx, laterBatch)
	defer blr.Close()

	var requestUpdates []string
	var sendUpdates []models.Data
	for i := 0; i < len(syncData); i++ {
		rowsE, err := ber.Query()
		if err != nil {
			return nil, nil, err
		}
		for rowsE.Next() {
			var id string
			err = rowsE.Scan(&id)
			if err != nil {
				return nil, nil, err
			}
			requestUpdates = append(requestUpdates, id)
		}
		rowsL, err := blr.Query()
		if err != nil {
			return nil, nil, err
		}
		if rowsL.Err() != nil {
			return nil, nil, rowsL.Err()
		}
		for rowsL.Next() {
			var data models.Data
			var name, dataType sql.NullString
			err = rowsL.Scan(&data.ID, &name, &dataType, &data.Content, &data.UpdatedAt, &data.Deleted)
			if err != nil {
				return nil, nil, err
			}
			data.Name = name.String
			data.Type = dataType.String
			sendUpdates = append(sendUpdates, data)
		}
	}
	return sendUpdates, requestUpdates, nil
}
