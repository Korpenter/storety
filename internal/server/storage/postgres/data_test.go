package postgres

import (
	"context"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestDB_CreateData(t *testing.T) {
	tests := []struct {
		name    string
		data    *models.Data
		userID  uuid.UUID
		resIns  pgconn.CommandTag
		resUpd  pgconn.CommandTag
		wantErr error
	}{
		{
			name:   "Create data successfully",
			resIns: pgxmock.NewResult("INSERT", 1),
			data: &models.Data{
				ID:        uuid.New(),
				Name:      "data",
				Type:      "binary",
				Content:   []byte{123},
				UpdatedAt: time.Now().UTC(),
				Deleted:   false,
			},
			userID:  uuid.New(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO data`)).
				WithArgs(tt.data.ID, tt.userID, tt.data.Name, tt.data.Type, tt.data.Content, tt.data.UpdatedAt, tt.data.Deleted).
				WillReturnResult(tt.resIns)
			mock.ExpectCommit()

			db := &DB{conn: mock}

			err = db.CreateData(context.Background(), tt.userID, tt.data)
			assert.ErrorIs(t, err, tt.wantErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetDataContentByName(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		rows     *pgxmock.Rows
		want     []byte
		wantType string
		dataName string
		wantErr  error
	}{
		{
			name: "Get existing data",
			rows: pgxmock.NewRows([]string{"content", "type"}).
				AddRow([]byte("content"), "Text"),
			want:     []byte("content"),
			wantType: "Text",
			dataName: "dataName",
			wantErr:  nil,
		},
		{
			name:     "Get with non-existent data",
			rows:     pgxmock.NewRows([]string{"content", "type"}),
			want:     nil,
			dataName: "dataName",
			wantErr:  constants.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			mock.ExpectQuery("SELECT content, type FROM data").
				WithArgs(tt.dataName, userID).WillReturnRows(tt.rows)
			db := &DB{conn: mock}
			content, contentType, err := db.GetDataContentByName(context.Background(), userID, tt.dataName)
			assert.Equal(t, tt.want, content)
			assert.Equal(t, tt.wantType, contentType)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteDataByName(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		resIns   pgconn.CommandTag
		resUpd   pgconn.CommandTag
		dataName string
		wantErr  error
	}{
		{
			name:     "Delete data successfully",
			resIns:   pgxmock.NewResult("INSERT", 1),
			resUpd:   pgxmock.NewResult("UPDATE", 1),
			dataName: "dataName",
			wantErr:  nil,
		},
		{
			name:     "Delete with non-existent data",
			resIns:   pgxmock.NewResult("DELETE", 0),
			dataName: "nonExistentName",
			wantErr:  constants.ErrDeleteData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE data SET`)).
				WithArgs(tt.dataName, userID).WillReturnResult(tt.resIns)
			if tt.wantErr == nil {
				mock.ExpectCommit()
			}
			db := &DB{conn: mock}
			err = db.DeleteDataByName(context.Background(), userID, tt.dataName)
			assert.ErrorIs(t, err, tt.wantErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDB_GetAllDataInfo(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name    string
		rows    *pgxmock.Rows
		want    []models.DataInfo
		userID  uuid.UUID
		wantErr error
	}{
		{
			name: "Get existing data info",
			rows: pgxmock.NewRows([]string{"name", "type"}).
				AddRow("dataName", "binary"),
			want:    []models.DataInfo{{Name: "dataName", Type: "binary"}},
			userID:  userID,
			wantErr: nil,
		},
		{
			name:    "Get with non-existent data info",
			rows:    pgxmock.NewRows([]string{"name", "type"}),
			want:    nil,
			userID:  uuid.New(),
			wantErr: constants.ErrNoData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			mock.ExpectQuery("SELECT name, type FROM data").
				WithArgs(tt.userID).WillReturnRows(tt.rows)
			db := &DB{conn: mock}
			list, err := db.GetAllDataInfo(context.Background(), tt.userID)
			assert.EqualValues(t, tt.want, list)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
