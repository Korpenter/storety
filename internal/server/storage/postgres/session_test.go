package postgres

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// CreateUser test
func TestDB_CreateSession(t *testing.T) {
	id, err := uuid.NewRandom()
	assert.NoError(t, err)

	tests := []struct {
		name    string
		rows    *pgxmock.Rows
		wantErr error
	}{
		{
			name: "Create session successfully",
			rows: pgxmock.NewRows([]string{"id", "user_id", "auth_token", "refresh_token"}).
				AddRow(id.String(), id.String(), "authToken", "refreshToken"),
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
			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions`)).
				WithArgs(id, id, "authToken", "refreshToken").WillReturnResult(pgxmock.NewResult("INSERT", 1))
			mock.ExpectCommit()

			db := &DB{conn: mock}
			session := &models.Session{
				ID:           id,
				UserID:       id,
				AuthToken:    "authToken",
				RefreshToken: "refreshToken",
			}
			err = db.CreateSession(context.Background(), session, nil)
			assert.ErrorIs(t, err, tt.wantErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDB_GetSession(t *testing.T) {
	userID, err := uuid.NewRandom()
	assert.NoError(t, err)
	sessionID, err := uuid.NewRandom()
	assert.NoError(t, err)
	tests := []struct {
		name         string
		rows         *pgxmock.Rows
		wantID       uuid.UUID
		refreshToken string
		wantErr      error
	}{
		{
			name: "Get with existing session",
			rows: pgxmock.NewRows([]string{"user_id"}).
				AddRow(userID.String()),
			wantID:       userID,
			refreshToken: "refreshToken",
			wantErr:      nil,
		},
		{
			name:         "Get with non-existent session",
			rows:         pgxmock.NewRows([]string{"user_id"}),
			wantID:       uuid.Nil,
			refreshToken: "",
			wantErr:      storage.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			mock.ExpectQuery("SELECT user_id FROM sessions").
				WithArgs(sessionID, "refreshToken").WillReturnRows(tt.rows)
			db := &DB{conn: mock}
			uID, err := db.GetSession(context.Background(), sessionID, "refreshToken")
			assert.EqualValues(t, tt.wantID, uID)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
